package validator

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type valid struct {
	req            *http.Request
	ruleRowList    []ruleRow // 规则列表
	ruleAsDataMap  map[string]*ruleAsData
	ruleAsDataList []ruleAsData
	errors         []error // 错误列表
}

// 验证json数据
func (v *valid) ValidJson(args ...interface{}) (va *valid) {
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
	if err := v.parseJsonData(data); err != nil {
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

// 拆分规则和数据成一个数据对应一个规则
func (v *valid) splitRuleAsData(row ruleRow, data interface{}, fullPk string) {
	if v.ruleAsDataMap == nil {
		v.ruleAsDataMap = map[string]*ruleAsData{}
	}
	v.ruleAsDataMap[fullPk] = &ruleAsData{
		methods: row.methods,
		notes:   row.notes,
		data:    data,
	}
	v.ruleAsDataList = append(v.ruleAsDataList, ruleAsData{
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
				v.splitRuleAsData(childRow, data, getFullKey(fullPk, childRow.field))
			}
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
				validData:     row.data,
				ruleAsDataMap: v.ruleAsDataMap,
			}
			var fn methodFunc
			switch m.method.(type) {
			case string:
				if m.method.(string) == "errors" {
					isErrors = true
					continue
				}
				var fnInterface interface{}
				var ok bool
				if fnInterface, ok = methodPool.Load(m.method.(string)); !ok {
					es = append(es, fmt.Errorf("规则 %s 不存在", m.method))
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
			if err := fn(d, m.args...); err != nil {
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
