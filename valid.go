package validator

import (
	"encoding/json"
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
}

// 验证json数据
func (v *valid) ValidJson(args ...interface{}) error {
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
		return errors.New("验证规则*rules和值*data至少传一个")
	}
	if err := v.parseRules(rs); err != nil {
		return err
	}
	if err := v.parseData(data); err != nil {
		return err
	}
	dataValue := reflect.ValueOf(data).Elem()
	newData := dataValue.Interface()
	for _, row := range v.ruleRowList {
		v.splitRuleAsData(row, newData, "")
	}
	if err := v.validRule(&newData); err != nil {
		return err
	}
	dataValue.Set(reflect.ValueOf(newData))
	return nil
}

// 解析数据
func (v *valid) parseData(data interface{}) error {
	body := readBody(v.req)
	if err := json.Unmarshal(body, data); err != nil {
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
func (v *valid) validRule(data *interface{}) error {
	for _, row := range v.ruleAsDataList {
		// 验证数据
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
				var fnInterface interface{}
				var ok bool
				if fnInterface, ok = methodPool.Load(m.method.(string)); !ok {
					return fmt.Errorf("规则 %s 不存在", m.method)
				}
				if fn, ok = fnInterface.(methodFunc); !ok {
					return errors.New("规则没有注入")
				}

			case methodFunc:
				fn = m.method.(methodFunc)
			default:
				return errors.New("未知错误")
			}
			if err := fn(d, m.args...); err != nil {
				return err
			}
		}
	}
	return nil
}
