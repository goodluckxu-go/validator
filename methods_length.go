package validator

import (
	"fmt"
	"unicode/utf8"
)

func (m *methods) Len(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	dataLen, ok := args[0].(int)
	if !ok {
		return fmt.Errorf("验证规则错误: 参数数量必须是integer类型")
	}
	rsErr := getMessageError(lang.Len, d.message, d.GetNotes(), dataLen)
	// 如果是文件则验证
	f := d.handle.fileMap[d.fullField]
	if f != nil {
		if int(f.Size/1024) == dataLen {
			return nil
		}
		return rsErr
	}
	validData := d.GetValidData()
	switch validData.(type) {
	case string:
		if utf8.RuneCountInString(validData.(string)) == dataLen {
			return nil
		}
	case []interface{}:
		if len(validData.([]interface{})) == dataLen {
			return nil
		}
	}
	return rsErr
}

func (m *methods) Min(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	dataLen, ok := args[0].(int)
	if !ok {
		return fmt.Errorf("验证规则错误: 参数数量必须是integer类型")
	}
	rsErr := getMessageError(lang.Min, d.message, d.GetNotes(), dataLen)
	// 如果是文件则验证
	f := d.handle.fileMap[d.fullField]
	if f != nil {
		if int(f.Size/1024) >= dataLen {
			return nil
		}
		return rsErr
	}
	validData := d.GetValidData()
	switch validData.(type) {
	case string:
		if utf8.RuneCountInString(validData.(string)) >= dataLen {
			return nil
		}
	case []interface{}:
		if len(validData.([]interface{})) >= dataLen {
			return nil
		}
	}
	return rsErr
}

func (m *methods) Max(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	dataLen, ok := args[0].(int)
	if !ok {
		return fmt.Errorf("验证规则错误: 参数数量必须是integer类型")
	}
	rsErr := getMessageError(lang.Max, d.message, d.GetNotes(), dataLen)
	// 如果是文件则验证
	f := d.handle.fileMap[d.fullField]
	if f != nil {
		if int(f.Size/1024) <= dataLen {
			return nil
		}
		return rsErr
	}
	validData := d.GetValidData()
	switch validData.(type) {
	case string:
		if utf8.RuneCountInString(validData.(string)) <= dataLen {
			return nil
		}
	case []interface{}:
		if len(validData.([]interface{})) <= dataLen {
			return nil
		}
	}
	return rsErr
}
