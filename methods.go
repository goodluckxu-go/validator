package validator

import "github.com/goodluckxu-go/validator/param"

type methods struct {
}

func (m *methods) Errors(d *Data, args ...*param.Param) error {
	return nil
}

func (m *methods) Required(d *Data, args ...*param.Param) error {
	rsErr := getMessageError(lang.Required, d.message, d.GetNotes())
	validData := d.GetValidData()
	if validData == nil {
		return rsErr
	}
	switch validData.(type) {
	case string:
		if validData.(string) == "" {
			return rsErr
		}
	case float64:
		if validData.(float64) == 0 {
			return rsErr
		}
	case bool:
		if validData.(bool) == false {
			return rsErr
		}
	case []interface{}:
		if len(validData.([]interface{})) == 0 {
			return rsErr
		}
	case map[string]interface{}:
		if len(validData.(map[string]interface{})) == 0 {
			return rsErr
		}
	}
	return nil
}

func (m *methods) ValidField(d *Data, args ...*param.Param) error {
	return nil
}
