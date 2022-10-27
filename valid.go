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

type valid struct {
	storage *storage // 存储仓库
	handle  *handle  // 处理数据
	Error   error
	Errors  []error
}

// Valid 验证数据
func (v *valid) Valid() (va *valid) {
	va = v.getInstance()
	var err error
	if va.storage.req != nil {
		err = va.parseRequest()
		if err != nil {
			va.Errors = []error{err}
			va.Error = err
			return
		}
	}
	if va.storage.rules != nil {
		var ruleRowList []ruleRow
		if ruleRowList, err = va.parseRules(va.storage.rules); err != nil {
			va.Errors = []error{err}
			va.Error = err
			return
		}
		dataValue := reflect.ValueOf(va.storage.data).Elem()
		newData := dataValue.Interface()
		va.handleRules(newData, ruleRowList)
		va.handleMessageOrNotes(va.storage.rules, va.storage.messages, newData)
		if errs := va.validRule(&newData); errs != nil {
			va.Errors = errs
			va.Error = errs[0]
			return
		}
		dataValue.Set(reflect.ValueOf(newData))
	}
	return
}

// 重新实例化
func (v *valid) getInstance() (va *valid) {
	va = new(valid)
	va.storage = v.storage
	va.handle = v.handle
	return
}

// parseRequest 解析请求数据
func (v *valid) parseRequest() error {
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
				_, fh, err := v.storage.req.FormFile(name)
				if err != nil {
					return err
				}
				v.handle.fileMap[name] = &file{
					Suffix: strings.TrimPrefix(filepath.Ext(fh.Filename), "."),
					Mime:   fh.Header.Get("Content-Type"),
					Name:   fh.Filename,
					Size:   fh.Size,
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
func (v *valid) parseRules(rules []Rule) ([]ruleRow, error) {
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
func (v *valid) handleRules(data interface{}, ruleRowList []ruleRow) {
	for _, row := range ruleRowList {
		v.splitRuleAsData(row, data, "")
	}
}

// 处理消息和注释
func (v *valid) handleMessageOrNotes(rules []Rule, messages []Message, data interface{}) {
	for _, r := range rules {
		v.splitMessages([3]string{r.Field, r.Notes, "node"}, data, "", 0)
	}
	for _, msg := range messages {
		v.splitMessages([3]string{msg[0], msg[1], "message"}, data, "", 0)
	}
	if v.handle.messageMap == nil {
		v.handle.messageMap = map[string]string{}
	}
	if v.handle.notesMap == nil {
		v.handle.notesMap = map[string]string{}
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
func (v *valid) splitRuleAsData(row ruleRow, data interface{}, fullPk string) {
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
func (v *valid) splitMessages(message [3]string, data interface{}, fullPk string, ln int) {
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
func (v *valid) validRule(data *interface{}) (es []error) {
	for _, fullPk := range v.handle.ruleIndex {
		row := v.handle.ruleData[fullPk]
		// 验证数据
		isErrors := false // 是否返回多个错误信息
		for _, m := range row.methods {
			if row.notes == "" {
				row.notes = fullPk
			}
			d := &Data{
				data:          data,
				notes:         row.notes,
				fullField:     fullPk,
				pk:            row.pk,
				validData:     &row.data,
				ruleAsDataMap: v.handle.ruleData,
				messageMap:    v.handle.messageMap,
				notesMap:      v.handle.notesMap,
			}
			var fn methodFunc
			var me string
			switch m.method.(type) {
			case string:
				me = m.method.(string)
				if me == "errors" {
					isErrors = true
					continue
				}
				d.message = v.handle.messageMap[fullPk+"."+me]
				var fnInterface interface{}
				var ok bool
				if fnInterface, ok = methodPool.Load(me); !ok {
					es = append(es, fmt.Errorf("规则 %s 不存在", me))
					return
				}
				if fn, ok = fnInterface.(methodFunc); !ok {
					es = append(es, errors.New("规则没有注入"))
					return
				}

			case methodFunc:
				fn = m.method.(methodFunc)
			default:
				es = append(es, errors.New("未知错误"))
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
				es = append(es, err)
			}
			if !isErrors && len(es) > 0 {
				return
			}
		}
		if isErrors && len(es) > 0 {
			return
		}
	}
	return nil
}
