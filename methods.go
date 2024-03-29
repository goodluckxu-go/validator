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
	// 先验证是否为文件
	if _, err := d.getFile(); err == nil {
		return nil
	}
	validData := d.GetValidData()
	if validData == nil {
		return validError(lang.Required, d.getMessage(), langArg{
			notes: d.GetNotes(),
		})
	}
	switch validData.(type) {
	case string:
		if validData.(string) == "" && fType != "string" {
			return validError(lang.Required, d.getMessage(), langArg{
				notes: d.GetNotes(),
			})
		}
	case float64:
		if validData.(float64) == 0 && fType != "number" {
			return validError(lang.Required, d.getMessage(), langArg{
				notes: d.GetNotes(),
			})
		}
	case bool:
		if validData.(bool) == false && fType != "bool" {
			return validError(lang.Required, d.getMessage(), langArg{
				notes: d.GetNotes(),
			})
		}
	case []interface{}:
		if len(validData.([]interface{})) == 0 && fType != "array" {
			return validError(lang.Required, d.getMessage(), langArg{
				notes: d.GetNotes(),
			})
		}
	case map[string]interface{}:
		if len(validData.(map[string]interface{})) == 0 && fType != "map" {
			return validError(lang.Required, d.getMessage(), langArg{
				notes: d.GetNotes(),
			})
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
		return d.JumpChild()
	}
	return nil
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
		return d.JumpChild()
	}
	return nil
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
			if toFloat64(validData) == toFloat64(child) {
				return nil
			}
		}
	}
	return validError(lang.In, d.getMessage(), langArg{
		notes: d.GetNotes(),
		array: args[0],
	})
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
				return validError(lang.NotIn, d.getMessage(), langArg{
					notes: d.GetNotes(),
					array: args[0],
				})
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			if toFloat64(validData) == toFloat64(child) {
				return validError(lang.NotIn, d.getMessage(), langArg{
					notes: d.GetNotes(),
					array: args[0],
				})
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
	return validError(lang.Email, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
}

func (m *methods) Phone(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	reg := regexp.MustCompile(`^1[0-9]{10}$`)
	if reg.MatchString(fmt.Sprintf("%v", d.GetValidData())) {
		return nil
	}
	return validError(lang.Phone, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
}

func (m *methods) Unique(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	pathList := strings.Split(d.path, ".")
	n := len(pathList)
	newList := make([]string, n)
	isExists := false
	for i := n - 1; i >= 0; i-- {
		tmp := pathList[i]
		if isExists {
			newList[i] = tmp
		} else if isSlice(tmp) {
			isExists = true
			newList[i] = "*"
		} else {
			newList[i] = tmp
		}
	}
	newPath := strings.Join(newList, ".")
	dataList := d.GetData(newPath)
	validData := d.GetValidData()
	for _, v := range dataList {
		if v.Path == d.path {
			continue
		}
		if validData == v.Data {
			return validError(lang.Unique, d.getMessage(), langArg{
				notes: d.GetNotes(),
			})
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
	return validError(lang.Regexp, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
}

func (m *methods) NotRegexp(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	reg := regexp.MustCompile(fmt.Sprintf("%v", args[0]))
	if !reg.MatchString(fmt.Sprintf("%v", d.GetValidData())) {
		return nil
	}
	return validError(lang.NotRegexp, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
}
