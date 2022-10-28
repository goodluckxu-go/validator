package validator

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Valid struct {
	storage storage // 存储仓库
	handle  handle  // 处理数据
	Error   error
}

// Valid 验证数据
func (v *Valid) Valid() (va *Valid) {
	va = v.getInstance()
	var err error
	if va.storage.req != nil {
		err = va.parseRequest()
		if err != nil {
			va.Error = err
			return
		}
	}
	if va.storage.rules != nil {
		var ruleRowList []ruleRow
		if ruleRowList, err = va.parseRules(va.storage.rules); err != nil {
			va.Error = err
			return
		}
		dataValue := reflect.ValueOf(va.storage.data).Elem()
		newData := dataValue.Interface()
		isSet := false
		switch newData.(type) {
		case map[string]interface{}, []interface{}:
			isSet = true
		default:
			by, byErr := json.Marshal(newData)
			if byErr != nil {
				va.Error = byErr
				return
			}
			if err = json.Unmarshal(by, &newData); err != nil {
				va.Error = err
				return
			}
		}
		va.handleRules(newData, ruleRowList)
		va.handleMessageOrNotes(va.storage.rules, va.storage.messages, newData)
		if errs := va.validRule(&newData); errs != nil {
			va.Error = errs
			return
		}
		if isSet {
			dataValue.Set(reflect.ValueOf(newData))
		}
	}
	return
}

// 重新实例化
func (v *Valid) getInstance() (va *Valid) {
	va = new(Valid)
	va.storage = v.storage
	va.handle = v.handle
	return
}

// parseRequest 解析请求数据
func (v *Valid) parseRequest() error {
	if v.storage.req == nil {
		return nil
	}
	contentTypeList := strings.Split(v.storage.req.Header.Get("Content-Type"), ";")
	contentType := contentTypeList[0]
	switch contentType {
	case "application/json":
		body := readBody(v.storage.req)
		if err := json.Unmarshal(body, v.storage.data); err != nil {
			return err
		}
	case "application/xml":
		body := readBody(v.storage.req)
		if err := xml.Unmarshal(body, v.storage.data); err != nil {
			return err
		}
	case "application/x-www-form-urlencoded":
		body := readBody(v.storage.req)
		bodyList := strings.Split(string(body), "&")
		data := map[string]interface{}{}
		for _, oneBody := range bodyList {
			oneList := strings.Split(oneBody, "=")
			oneKey, _ := url.PathUnescape(oneList[0])
			oneVal, _ := url.PathUnescape(oneList[1])
			data[oneKey] = oneVal
		}
		reflect.ValueOf(v.storage.data).Elem().Set(reflect.ValueOf(data))
	case "multipart/form-data":
		if len(contentTypeList) < 2 {
			return errors.New("请传入 multipart/form-data 格式数据")
		}
		boundary := strings.Split(contentTypeList[1], "=")[1]
		body := readBody(v.storage.req)
		bodyList := strings.Split(string(body), boundary)
		bodyList = bodyList[1 : len(bodyList)-1]
		data := map[string]interface{}{}
		for _, oneBody := range bodyList {
			oneList := strings.Split(oneBody, "\r\n")
			explanList := strings.Split(oneList[1], ";")
			name := strings.Trim(strings.Split(explanList[1], "=")[1], "\"")
			if len(explanList) == 2 {
				// 常规数据
				data[name] = oneList[3]
			} else {
				// 文件数据
				if _, fh, err := v.storage.req.FormFile(name); err == nil {
					v.handle.fileMap[name] = &file{
						Suffix: strings.TrimPrefix(filepath.Ext(fh.Filename), "."),
						Mime:   fh.Header.Get("Content-Type"),
						Name:   fh.Filename,
						Size:   fh.Size,
					}
				}
			}
		}
		reflect.ValueOf(v.storage.data).Elem().Set(reflect.ValueOf(data))
	default:
		body := readBody(v.storage.req)
		reflect.ValueOf(v.storage.data).Elem().Set(reflect.ValueOf(string(body)))
	}
	return nil
}

// 解析规则
func (v *Valid) parseRules(rules []Rule) ([]ruleRow, error) {
	if rules == nil {
		return nil, nil
	}
	ruleList, err := disintegrateRules(rules)
	if err != nil {
		return nil, err
	}
	return assembleRuleRow(ruleList), nil
}

// 处理规则
func (v *Valid) handleRules(data interface{}, ruleRowList []ruleRow) {
	for _, row := range ruleRowList {
		v.splitRuleAsData(row, data, "")
	}
}

// 处理消息和注释
func (v *Valid) handleMessageOrNotes(rules []Rule, messages []Message, data interface{}) {
	for _, r := range rules {
		v.splitMessages([3]string{r.Field, r.Notes, "node"}, data, "", 0)
	}
	for _, msg := range messages {
		v.splitMessages([3]string{msg[0], msg[1], "message"}, data, "", 0)
	}
	for _, r := range v.handle.messages {
		if r[2] == "message" {
			v.handle.messageMap[r[0]] = r[1]
		} else if r[2] == "node" {
			if r[1] == "" {
				r[1] = r[0]
			}
			v.handle.notesMap[r[0]] = r[1]
		}
	}
}

// 拆分规则和数据成一个数据对应一个规则
func (v *Valid) splitRuleAsData(row ruleRow, data interface{}, fullPk string) {
	pk := strings.TrimPrefix(row.pk, "root")
	pk = strings.TrimPrefix(pk, ".")
	if v.handle.ruleData == nil {
		v.handle.ruleData = map[string]*ruleAsData{}
	}
	v.handle.ruleData[fullPk] = &ruleAsData{
		pk:      pk,
		methods: row.methods,
		notes:   row.notes,
		data:    data,
	}
	v.handle.ruleIndex = append(v.handle.ruleIndex, fullPk)
	for _, childRow := range row.children {
		if childRow.field == "*" {
			dataList, _ := data.([]interface{})
			for index, childData := range dataList {
				v.splitRuleAsData(childRow, childData, getFullKey(fullPk, index))
			}
			if len(dataList) == 0 {
				v.splitRuleAsData(childRow, nil, getFullKey(fullPk, 0))
			}
		} else {
			if dataMap, ok := data.(map[string]interface{}); ok {
				v.splitRuleAsData(childRow, dataMap[childRow.field], getFullKey(fullPk, childRow.field))
			} else {
				dataList, _ := data.([]interface{})
				index, _ := strconv.Atoi(childRow.field)
				var childData interface{}
				if len(dataList) > index {
					childData = dataList[index]
				}
				v.splitRuleAsData(childRow, childData, getFullKey(fullPk, childRow.field))
			}
		}
	}
}

// 拆分规则和消息成一个数据对应一个规则
func (v *Valid) splitMessages(message [3]string, data interface{}, fullPk string, ln int) {
	msgKeyList := strings.Split(message[0], ".")
	fullPkList := strings.Split(fullPk, ".")
	if ln > 0 && len(fullPkList) == ln {
		v.handle.messages = append(v.handle.messages, [3]string{fullPk, message[1], message[2]})
	}
	firstField := msgKeyList[0]
	if firstField == "" {
		return
	}
	if ln == 0 {
		ln = len(msgKeyList)
	}
	otherField := strings.Join(msgKeyList[1:], ".")
	message[0] = otherField
	if firstField == "*" {
		dataList, _ := data.([]interface{})
		v.splitMessages(message, dataList, getFullKey(fullPk, "*"), ln)
		for index, childData := range dataList {
			v.splitMessages(message, childData, getFullKey(fullPk, index), ln)
		}
		if len(dataList) == 0 {
			v.splitMessages(message, nil, getFullKey(fullPk, 0), ln)
		}
	} else {
		if dataMap, ok := data.(map[string]interface{}); ok {
			v.splitMessages(message, dataMap[firstField], getFullKey(fullPk, firstField), ln)
		} else {
			dataList, _ := data.([]interface{})
			index, _ := strconv.Atoi(firstField)
			var childData interface{}
			if len(dataList) > index {
				childData = dataList[index]
			}
			v.splitMessages(message, childData, getFullKey(fullPk, firstField), ln)
		}
	}
}

// 验证规则
func (v *Valid) validRule(data *interface{}) (errs error) {
	for _, fullPk := range v.handle.ruleIndex {
		row := v.handle.ruleData[fullPk]
		fullPkList := strings.Split(fullPk, ".")
		if len(fullPkList) > 1 {
			parentFullKey := strings.Join(fullPkList[0:len(fullPkList)-1], ".")
			parentRow := v.handle.ruleData[parentFullKey]
			if parentRow != nil && !parentRow.isValid {
				continue
			}
		}
		// 验证数据
		validNum := 0
		for _, m := range row.methods {
			validNum++
			if row.notes == "" {
				row.notes = fullPk
			}
			d := &Data{
				data:      data,
				notes:     row.notes,
				fullField: fullPk,
				pk:        row.pk,
				validData: &row.data,
				handle:    v.handle,
			}
			var fn methodFunc
			var me string
			switch m.method.(type) {
			case string:
				me = m.method.(string)
				d.message = v.handle.messageMap[fullPk+"."+me]
				var fnInterface interface{}
				var ok bool
				if fnInterface, ok = methodPool.Load(me); !ok {
					errs = fmt.Errorf("规则 %s 不存在", me)
					return
				}
				if fn, ok = fnInterface.(methodFunc); !ok {
					errs = errors.New("规则没有注入")
					return
				}

			case methodFunc:
				fn = m.method.(methodFunc)
			default:
				errs = errors.New("未知错误")
				return
			}
			err := fn(d, m.args...)
			// 是否验证下面数据
			if inArray(me, []string{"valid_condition", "nullable"}) {
				if err == nil {
					break
				}
				if err.Error() == "" {
					continue
				}
			}
			if err != nil {
				errs = err
				return
			}
		}
		if len(row.methods) == validNum {
			v.handle.ruleData[fullPk].isValid = true
		}
	}
	return nil
}
