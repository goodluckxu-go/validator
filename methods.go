package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type methods struct {
}

func (m *methods) Required(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 1, []string{"string", "number", "bool", "array", "map"}); err != nil {
		return err
	}
	var fType string
	if len(args) > 0 {
		fType, _ = args[0].(string)
	}
	rsErr := getMessageError(lang.Required, d.message, d.GetNotes())
	// 先验证是否为文件
	if d.handle.fileMap[d.fullField] != nil {
		return nil
	}
	validData := d.GetValidData()
	if validData == nil {
		return rsErr
	}
	switch validData.(type) {
	case string:
		if validData.(string) == "" && fType != "string" {
			return rsErr
		}
	case float64:
		if validData.(float64) == 0 && fType != "number" {
			return rsErr
		}
	case bool:
		if validData.(bool) == false && fType != "bool" {
			return rsErr
		}
	case []interface{}:
		if len(validData.([]interface{})) == 0 && fType != "array" {
			return rsErr
		}
	case map[string]interface{}:
		if len(validData.(map[string]interface{})) == 0 && fType != "map" {
			return rsErr
		}
	}
	return nil
}

func (m *methods) ValidCondition(d *Data, args ...interface{}) error {
	if err := validArgs(args, 2, -1); err != nil {
		return err
	}
	bl, err := formulaCompare(d, args...)
	if err != nil {
		return err
	}
	if !bl {
		return d.JumpValid()
	}
	return d.NextValid()
}

func (m *methods) Nullable(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	validData := d.GetValidData()
	isNull := false
	if validData == nil {
		isNull = true
	} else {
		switch validData.(type) {
		case string:
			if validData.(string) == "" {
				isNull = true
			}
		case float64:
			if validData.(float64) == 0 {
				isNull = true
			}
		case bool:
			if validData.(bool) == false {
				isNull = true
			}
		case []interface{}:
			if len(validData.([]interface{})) == 0 {
				isNull = true
			}
		case map[string]interface{}:
			if len(validData.(map[string]interface{})) == 0 {
				isNull = true
			}
		}
	}
	if isNull {
		return d.JumpValid()
	}
	return d.NextValid()
}

func (m *methods) In(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	arrValue := reflect.ValueOf(args[0])
	if arrValue.Kind() != reflect.Slice && arrValue.Kind() != reflect.Array {
		return fmt.Errorf("验证规则错误: 参数类型必须是数组或切片")
	}
	validData := d.GetValidData()
	arrLen := arrValue.Len()
	for i := 0; i < arrLen; i++ {
		child := arrValue.Index(i).Interface()
		switch child.(type) {
		case string:
			if fmt.Sprintf("%v", validData) == child.(string) {
				return nil
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			if interfaceToFloat64(validData) == interfaceToFloat64(child) {
				return nil
			}
		}
	}
	return getMessageError(lang.In, d.message, d.GetNotes(), fmt.Sprintf("%v", args[0]))
}

func (m *methods) NotIn(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	arrValue := reflect.ValueOf(args[0])
	if arrValue.Kind() != reflect.Slice && arrValue.Kind() != reflect.Array {
		return fmt.Errorf("验证规则错误: 参数类型必须是数组或切片")
	}
	validData := d.GetValidData()
	arrLen := arrValue.Len()
	for i := 0; i < arrLen; i++ {
		child := arrValue.Index(i).Interface()
		switch child.(type) {
		case string:
			if fmt.Sprintf("%v", validData) == child.(string) {
				return getMessageError(lang.NotIn, d.message, d.GetNotes(), fmt.Sprintf("%v", args[0]))
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			if interfaceToFloat64(validData) == interfaceToFloat64(child) {
				return getMessageError(lang.NotIn, d.message, d.GetNotes(), fmt.Sprintf("%v", args[0]))
			}
		}
	}
	return nil
}

func (m *methods) Email(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	reg := regexp.MustCompile(`^[A-Za-z0-9\\u4e00-\\u9fa5]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	if reg.MatchString(fmt.Sprintf("%v", d.GetValidData())) {
		return nil
	}
	return getMessageError(lang.Email, d.message, d.GetNotes())
}

func (m *methods) Phone(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	reg := regexp.MustCompile(`^1[0-9]{10}$`)
	if reg.MatchString(fmt.Sprintf("%v", d.GetValidData())) {
		return nil
	}
	return getMessageError(lang.Phone, d.message, d.GetNotes())
}

func (m *methods) Unique(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	pkList := strings.Split(d.pk, ".")
	lastStarIndex := 0
	for k, v := range pkList {
		if v == "*" {
			lastStarIndex = k
		}
	}
	fullList := strings.Split(d.fullField, ".")
	newList := append(fullList[0:lastStarIndex], pkList[lastStarIndex:]...)
	var data interface{}
	for k, v := range d.handle.ruleData {
		kList := strings.Split(k, ".")
		isEq := true
		for i, n := range newList {
			if len(kList) <= i {
				isEq = false
				break
			}
			if n == "*" || kList[i] == n {
				continue
			}
			isEq = false
		}
		if isEq {
			if data == nil {
				data = v.data
			} else {
				if data == v.data {
					return getMessageError(lang.Unique, d.message, d.GetNotes())
				}
			}
		}
	}
	return nil
}

func (m *methods) Regexp(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	reg := regexp.MustCompile(fmt.Sprintf("%v", args[0]))
	if reg.MatchString(fmt.Sprintf("%v", d.GetValidData())) {
		return nil
	}
	return getMessageError(lang.Regexp, d.message, d.GetNotes())
}

func (m *methods) NotRegexp(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	reg := regexp.MustCompile(fmt.Sprintf("%v", args[0]))
	if !reg.MatchString(fmt.Sprintf("%v", d.GetValidData())) {
		return nil
	}
	return getMessageError(lang.NotRegexp, d.message, d.GetNotes())
}
