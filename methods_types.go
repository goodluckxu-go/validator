package validator

import "strconv"

func (m *methods) Array(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	if _, ok := validData.([]interface{}); ok {
		return nil
	}
	return validError(lang.Array, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
}

func (m *methods) Map(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	if _, ok := validData.(map[string]interface{}); ok {
		return nil
	}
	return validError(lang.Map, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
}
func (m *methods) String(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	if _, ok := validData.(string); ok {
		return nil
	}
	return validError(lang.String, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
}
func (m *methods) Number(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
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
	return validError(lang.Number, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
}

func (m *methods) Integer(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	switch validData.(type) {
	case float64:
		validDataInt := int64(validData.(float64))
		if validData.(float64) == float64(validDataInt) {
			d.setValidData(validDataInt)
			return nil
		}
	case string:
		if number, err := strconv.ParseFloat(validData.(string), 64); err == nil {
			if number == float64(int64(number)) {
				d.setValidData(int64(number))
				return nil
			}
		}
	}
	return validError(lang.Integer, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
}

func (m *methods) Bool(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 0); err != nil {
		return err
	}
	validData := d.GetValidData()
	if validData == nil {
		return nil
	}
	if _, ok := validData.(bool); ok {
		return nil
	}
	return validError(lang.Bool, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
}

func (m *methods) Date(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 1); err != nil {
		return err
	}
	rsErr := validError(lang.Date, d.getMessage(), langArg{
		notes: d.GetNotes(),
	})
	dateString, ok := d.GetValidData().(string)
	if !ok {
		return rsErr
	}
	if len(args) == 0 {
		if _, err := timeParse(dateString); err != nil {
			return rsErr
		}
		return nil
	}
	formatString, _ := args[0].(string)
	if err := validDate(dateString, formatString); err != nil {
		return rsErr
	}
	return nil
}

func (m *methods) File(d *Data, args ...interface{}) error {
	if err := validArgs(args, 0, 1); err != nil {
		return err
	}
	if _, err := d.getFile(); err != nil {
		return validError(lang.File, d.getMessage(), langArg{
			notes: d.GetNotes(),
		})
	}
	return nil
}
