package validator

import (
	"errors"
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
		return nil
	}
	return errors.New("")
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
		return nil
	}
	return errors.New("")
}
