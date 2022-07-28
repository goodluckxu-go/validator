package validator

import (
	"fmt"
)

type methods struct {
}

func (m *methods) Required(d *Data, args ...interface{}) error {
	rs := getNotes(lang.Required, d.GetNotes())
	validData := d.GetValidData()
	if validData == nil {
		return rs
	}
	switch validData.(type) {
	case string:
		if validData.(string) == "" {
			return rs
		}
	case float64:
		if validData.(float64) == 0 {
			return rs
		}
	case bool:
		if validData.(bool) == false {
			return rs
		}
	case []interface{}:
		if len(validData.([]interface{})) == 0 {
			return rs
		}
	case map[string]interface{}:
		if len(validData.(map[string]interface{})) == 0 {
			return rs
		}
	}
	return nil
}

func (m *methods) ValidField(d *Data, args ...interface{}) error {
	fmt.Println(d.validData, d.GetValidData(), args)
	return nil
}
