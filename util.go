package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// 读取body内容
func readBody(r *http.Request) []byte {
	var bodyBytes []byte // 我们需要的body内容
	// 从原有Request.Body读取
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return bodyBytes
	}
	// 新建缓冲区并替换原有Request.body
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

func disintegrateRules(args ...interface{}) ([]interface{}, error) {
	var rs []Rule
	var inList []interface{}
	extendsMap := map[string]interface{}{}
	for _, arg := range args {
		switch arg.(type) {
		case []Rule:
			rs = arg.([]Rule)
		case []interface{}:
			inList = arg.([]interface{})
		case map[string]interface{}:
			extendsMap = arg.(map[string]interface{})
		}
	}
	var list []interface{}
	if rs != nil {
		list = append(list, map[string]interface{}{
			"field": "root",
		})
		for _, v := range rs {
			newFiled := strings.Join([]string{
				"root",
				v.Field,
			}, ".")
			list = append(list, map[string]interface{}{
				"field":     newFiled,
				"pk":        "",
				"parent_pk": "",
			})
			extendsMap[newFiled] = map[string]interface{}{
				"methods": v.Methods,
				"notes":   v.Notes,
			}
		}
	} else if len(inList) > 0 {
		list = inList
	} else {
		return list, nil
	}
	// 如果传入数据为非json数据
	if len(rs) == 1 && rs[0].Field == "" {
		onlyOne := rs[0]
		return []interface{}{
			map[string]interface{}{
				"pk":        "root",
				"parent_pk": "",
			},
			map[string]interface{}{
				"pk":        "root.",
				"parent_pk": "root",
				"field":     "",
				"methods":   onlyOne.Methods,
				"notes":     onlyOne.Notes,
			},
		}, nil
	}
	var newList, otherList []interface{}
	parentMap := map[string]string{}
	var tmpList []string
	for _, v := range list {
		vMap, _ := v.(map[string]interface{})
		field, _ := vMap["field"].(string)
		fieldList := strings.Split(field, ".")
		firstField := fieldList[0]
		otherField := strings.Join(fieldList[1:], ".")
		parentPk, _ := vMap["parent_pk"].(string)
		var pkList []string
		if parentPk != "" {
			pkList = append(pkList, parentPk)
		}
		pkList = append(pkList, firstField)
		pk := strings.Join(pkList, ".")
		parentPkList := strings.Split(parentPk, ".")
		var completePkList []string
		completePkList = append(completePkList, pk)
		if otherField != "" {
			completePkList = append(completePkList, otherField)
		}
		completePk := strings.Join(completePkList, ".")
		completePk = strings.TrimPrefix(completePk, "root.")
		if parentMap[parentPk] == "" {
			parentMap[parentPk] = firstField
		}
		if (parentMap[parentPk] == "*" && firstField != "*") || (parentMap[parentPk] != "*" && firstField == "*") {
			tmpList = append(tmpList, "["+completePk+"]")
			return nil, fmt.Errorf("%s冲突，map和slice不能并存", strings.Join(tmpList, ","))
		}
		tmpList = []string{"[" + completePk + "]"}
		parentMap[parentPk] = firstField
		if !isMapInSliceMap(newList, map[string]interface{}{
			"field":     firstField,
			"parent_pk": parentPk,
		}) {
			extends, _ := extendsMap[pk].(map[string]interface{})
			newList = append(newList, map[string]interface{}{
				"pk":        pk,
				"parent_pk": parentPk,
				"field":     firstField,
				"methods":   extends["methods"],
				"notes":     extends["notes"],
			})
		}
		if !isMapInSliceMap(otherList, map[string]interface{}{
			"field":     otherField,
			"parent_pk": parentPk,
		}) && otherField != "" {
			parentPkList = []string{}
			if parentPk != "" {
				parentPkList = append(parentPkList, parentPk)
			}
			parentPkList = append(parentPkList, firstField)
			parentPk = strings.Join(parentPkList, ".")
			otherList = append(otherList, map[string]interface{}{
				"parent_pk": parentPk,
				"field":     otherField,
			})
		}
	}
	if len(otherList) > 0 {
		childList, err := disintegrateRules(otherList, extendsMap)
		if err != nil {
			return nil, err
		}
		newList = append(newList, childList...)
	}
	return newList, nil
}

func assembleRuleRow(list []interface{}, args ...interface{}) []ruleRow {
	parentPk := ""
	for _, arg := range args {
		switch arg.(type) {
		case string:
			parentPk = arg.(string)
		}
	}
	var newList []ruleRow
	for _, v := range list {
		vMap, _ := v.(map[string]interface{})
		if vMap["parent_pk"] == parentPk {
			pk, _ := vMap["pk"].(string)
			field, _ := vMap["field"].(string)
			m, _ := vMap["methods"].([]*method)
			var mList []methodData
			for _, mv := range m {
				mList = append(mList, mv.methods...)
			}
			notes, _ := vMap["notes"].(string)
			children := assembleRuleRow(list, pk)
			newList = append(newList, ruleRow{
				field:    field,
				pk:       pk,
				methods:  mList,
				notes:    notes,
				children: children,
			})
		}
	}
	return newList
}

func isMapInSliceMap(list []interface{}, where map[string]interface{}) bool {
	for _, v := range list {
		vMap, _ := v.(map[string]interface{})
		isEq := true
		for k, w := range where {
			if vMap[k] != w {
				isEq = false
			}
		}
		if len(where) == 0 {
			isEq = false
		}
		if isEq {
			return true
		}
	}
	return false
}

func getInstance(in interface{}) interface{} {
	var rs interface{}
	switch in.(type) {
	case *method:
		newIn, _ := in.(*method)
		out := new(method)
		out.methods = newIn.methods
		rs = out
	}
	return rs
}

func getFullKey(fullKey string, field interface{}) string {
	var fullKeyList []string
	if fullKey != "" {
		fullKeyList = append(fullKeyList, fullKey)
	}
	fullKeyList = append(fullKeyList, fmt.Sprintf("%v", field))
	return strings.Join(fullKeyList, ".")
}

func upData(data interface{}, key string, value interface{}) interface{} {
	if key == "" {
		return value
	}
	keyList := strings.Split(key, ".")
	firstKey := keyList[0]
	otherKey := strings.Join(keyList[1:], ".")
	switch data.(type) {
	case map[string]interface{}:
		dataMap, _ := data.(map[string]interface{})
		dataMap[firstKey] = upData(dataMap[firstKey], otherKey, value)
		data = dataMap
	case []interface{}:
		index, _ := strconv.Atoi(firstKey)
		dataList, _ := data.([]interface{})
		if len(dataList) > index {
			dataList[index] = upData(dataList[index], otherKey, value)
		}
		data = dataList
	}
	return data
}

func parseLang(langAddr string) {
	if langAddr == "" {
		return
	}
	var f *os.File
	var err error
	if f, err = os.Open(langAddr); err != nil {
		panic(err)
	}
	var by []byte
	if by, err = io.ReadAll(f); err != nil {
		panic(err)
	}
	langMap := map[string]interface{}{}
	if err = json.Unmarshal(by, &langMap); err != nil {
		panic(err)
	}
	langValue := reflect.ValueOf(&lang).Elem()
	langLen := langValue.NumField()
	for key, val := range langMap {
		for i := 0; i < langLen; i++ {
			field := langValue.Type().Field(i).Name
			if key == field {
				langValue.Field(i).Set(reflect.ValueOf(val))
				break
			}
		}
	}
}

func getNotes(langStr string, message string, args ...interface{}) error {
	if len(args) > 0 {
		langStr = strings.Replace(langStr, "${notes}", fmt.Sprintf("%v", args[0]), -1)
	}
	if message != "" {
		langStr = message
	}
	return errors.New(langStr)
}

func getCommonFullField(field, otherField string) (string, string) {
	fieldList := strings.Split(field, ".")
	otherFieldList := strings.Split(otherField, ".")
	var rsList []string
	i := 0
	for index, ov := range otherFieldList {
		if len(fieldList) <= index {
			break
		}
		fv := fieldList[index]
		if ov != "*" && ov != fv {
			break
		}
		rsList = append(rsList, fv)
		i++
	}
	return strings.Join(rsList, "."), strings.Join(otherFieldList[i:], ".")
}

func getData(data interface{}, key string) []interface{} {
	if key == "" {
		return []interface{}{data}
	}
	keyList := strings.Split(key, ".")
	firstKey := keyList[0]
	otherKey := strings.Join(keyList[1:], ".")
	switch data.(type) {
	case map[string]interface{}:
		dataMap, _ := data.(map[string]interface{})
		return getData(dataMap[firstKey], otherKey)
	case []interface{}:
		dataList, _ := data.([]interface{})
		var newDataList []interface{}
		for _, childData := range dataList {
			newDataList = append(newDataList, getData(childData, otherKey)...)
		}
		return newDataList
	}
	return []interface{}{data}
}

func getMessagesVal(messages []Message, row ruleAsData, me string) string {
	rs := ""
	fullPkList := strings.Split(row.fullPk, ".")
	for _, message := range messages {
		k := message[0]
		v := message[1]
		if k == "" {
			continue
		}
		kList := strings.Split(k, ".")
		if kList[len(kList)-1] != me {
			continue
		}
		kList = kList[0 : len(kList)-1]
		isSame := true
		for index, kv := range kList {
			if len(fullPkList) <= index {
				isSame = false
				break
			}
			fv := fullPkList[index]
			if kv != "*" && kv != fv {
				isSame = false
				break
			}
		}
		if isSame {
			rs = v
		}
	}
	return rs
}
