package validator

import (
	"strconv"
)

func (m *methods) Array(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	rsErr := getMessageError(lang.Array, d.message, d.GetNotes())
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	if _, ok := validData.([]interface{}); ok {
		return nil
	}
	return rsErr
}

func (m *methods) Map(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	rsErr := getMessageError(lang.Map, d.message, d.GetNotes())
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	if _, ok := validData.(map[string]interface{}); ok {
		return nil
	}
	return rsErr
}
func (m *methods) String(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	rsErr := getMessageError(lang.String, d.message, d.GetNotes())
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	if _, ok := validData.(string); ok {
		return nil
	}
	return rsErr
}
func (m *methods) Number(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	rsErr := getMessageError(lang.Number, d.message, d.GetNotes())
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	switch validData.(type) {
	case float64:
		return nil
	case string:
		if number, err := strconv.ParseFloat(validData.(string), 64); err == nil {
			d.setValidData(number)
			return nil
		}
	}
	return rsErr
}

func (m *methods) Integer(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	rsErr := getMessageError(lang.Integer, d.message, d.GetNotes())
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	switch validData.(type) {
	case float64:
		validDataInt := int(validData.(float64))
		if validData.(float64) == float64(validDataInt) {
			return nil
		}
	case string:
		if number, err := strconv.ParseFloat(validData.(string), 64); err == nil {
			if number == float64(int(number)) {
				d.setValidData(int(number))
				return nil
			}
		}
	}
	return rsErr
}

func (m *methods) Bool(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	rsErr := getMessageError(lang.Bool, d.message, d.GetNotes())
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	if _, ok := validData.(string); ok {
		return nil
	}
	return rsErr
}
