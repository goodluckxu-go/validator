package validator

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
	return nil
}
