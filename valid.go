package validator

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/goodluckxu-go/validator/param"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type valid struct {
	req            *http.Request
	ruleRowList    []ruleRow              // 规则列表
	notesMap       map[string]string      // 规则注释(notes)
	messageMap     map[string]string      // 规则注释(messages)
	ruleAsDataMap  map[string]*ruleAsData // 数据
	ruleAsDataList []ruleAsData
	messages       [][3]string
	errors         []error // 错误列表
}

// 验证json数据
func (v *valid) ValidJson(args ...interface{}) (va *valid) {
	va = new(valid)
	var rs []Rule
	var data interface{}
	var messages []Message
	for _, arg := range args {
		switch arg.(type) {
		case []Rule:
			rs = arg.([]Rule)
		case []Message:
			messages = arg.([]Message)
		case *map[string]interface{}, *[]interface{}, *interface{}:
			data = arg
		}
	}
	if data == nil && rs == nil {
		va.errors = append(va.errors, errors.New("验证规则*rules和值*data至少传一个"))
		return
	}
	if err := v.parseRules(rs); err != nil {
		va.errors = append(va.errors, err)
		return
	}
	if err := v.parseJsonData(data); err != nil {
		va.errors = append(va.errors, err)
		return
	}
	dataValue := reflect.ValueOf(data).Elem()
	newData := dataValue.Interface()
	v.handleRules(newData)
	v.handleMessageOrNotes(rs, messages, newData)
	if errs := v.validRule(&newData); errs != nil {
		va.errors = errs
		return
	}
	dataValue.Set(reflect.ValueOf(newData))
	return nil
}

// 验证xml数据
func (v *valid) ValidXml(args ...interface{}) (va *valid) {
	va = new(valid)
	var rs []Rule
	var data interface{}
	for _, arg := range args {
		switch arg.(type) {
		case []Rule:
			rs = arg.([]Rule)
		case *map[string]interface{}, *[]interface{}, *interface{}:
			data = arg
		}
	}
	if data == nil && rs == nil {
		va.errors = append(va.errors, errors.New("验证规则*rules和值*data至少传一个"))
		return
	}
	if err := v.parseRules(rs); err != nil {
		va.errors = append(va.errors, err)
		return
	}
	if err := v.parseXmlData(data); err != nil {
		va.errors = append(va.errors, err)
		return
	}
	dataValue := reflect.ValueOf(data).Elem()
	newData := dataValue.Interface()
	for _, row := range v.ruleRowList {
		v.splitRuleAsData(row, newData, "")
	}
	if errs := v.validRule(&newData); errs != nil {
		va.errors = errs
		return
	}
	dataValue.Set(reflect.ValueOf(newData))
	return nil
}

// 获取字符串错误信息列表
func (v *valid) Errors() (es []string) {
	for _, err := range v.errors {
		es = append(es, err.Error())
	}
	return
}

// 获取第一个错误信息
func (v *valid) Error() string {
	for _, err := range v.errors {
		return err.Error()
	}
	return ""
}

// 解析json数据
func (v *valid) parseJsonData(data interface{}) error {
	body := readBody(v.req)
	if err := json.Unmarshal(body, data); err != nil {
		return err
	}
	return nil
}

// 解析xml数据
func (v *valid) parseXmlData(data interface{}) error {
	body := readBody(v.req)
	if err := xml.Unmarshal(body, data); err != nil {
		return err
	}
	return nil
}

// 解析规则
func (v *valid) parseRules(rules []Rule) error {
	if rules == nil {
		return nil
	}
	ruleList, err := disintegrateRules(rules)
	if err != nil {
		return err
	}
	v.ruleRowList = assembleRuleRow(ruleList)
	return nil
}

// 处理规则
func (v *valid) handleRules(data interface{}) {
	for _, row := range v.ruleRowList {
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
	if v.messageMap == nil {
		v.messageMap = map[string]string{}
	}
	if v.notesMap == nil {
		v.notesMap = map[string]string{}
	}
	for _, r := range v.messages {
		if r[2] == "message" {
			v.messageMap[r[0]] = r[1]
		} else if r[2] == "node" {
			if r[1] == "" {
				r[1] = r[0]
			}
			v.notesMap[r[0]] = r[1]
		}
	}
}

// 拆分规则和数据成一个数据对应一个规则
func (v *valid) splitRuleAsData(row ruleRow, data interface{}, fullPk string) {
	pk := strings.TrimPrefix(row.pk, "root")
	pk = strings.TrimPrefix(pk, ".")
	if v.ruleAsDataMap == nil {
		v.ruleAsDataMap = map[string]*ruleAsData{}
	}
	v.ruleAsDataMap[fullPk] = &ruleAsData{
		pk:      pk,
		methods: row.methods,
		notes:   row.notes,
		data:    data,
	}
	v.ruleAsDataList = append(v.ruleAsDataList, ruleAsData{
		pk:      pk,
		fullPk:  fullPk,
		methods: row.methods,
		notes:   row.notes,
		data:    data,
	})
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
		v.messages = append(v.messages, [3]string{fullPk, message[1], message[2]})
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
	for _, row := range v.ruleAsDataList {
		// 验证数据
		isErrors := false // 是否返回多个错误信息
		for _, m := range row.methods {
			if row.notes == "" {
				row.notes = row.fullPk
			}
			d := &Data{
				data:          data,
				notes:         row.notes,
				fullField:     row.fullPk,
				validData:     &row.data,
				ruleAsDataMap: v.ruleAsDataMap,
				messageMap:    v.messageMap,
				notesMap:      v.notesMap,
			}
			var fn methodFunc
			switch m.method.(type) {
			case string:
				me := m.method.(string)
				if me == "errors" {
					isErrors = true
					continue
				}
				//d.message = getMessagesVal(v.messages, row.fullPk, me)
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
			var newArgs []*param.Param
			for _, arg := range m.args {
				if argParam, ok := arg.(*param.Param); ok {
					newArgs = append(newArgs, argParam)
				} else {
					newArgs = append(newArgs, param.SetValue(arg))
				}
			}
			if err := fn(d, newArgs...); err != nil {
				es = append(es, err)
			}
			if !isErrors {
				return
			}
		}
		if isErrors && len(es) > 0 {
			return
		}
	}
	return nil
}
