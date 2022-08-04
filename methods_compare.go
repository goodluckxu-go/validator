package validator

import (
	"test/validator/param"
	"test/validator/types"
)

func (m *methods) Eq(d *Data, args ...interface{}) error {
	var fType *param.Param
	var val interface{}
	for _, arg := range args {
		if argParam, ok := arg.(*param.Param); ok {
			switch argParam {
			case types.Field:
				fType = argParam
			default:
				if argParam.Value != nil {
					val = argParam.Value
				}
			}
		}
	}
	validData := d.GetValidData()
	validNotes := d.GetNotes()
	if fType == types.Field {
		valStr, _ := val.(string)
		for _, vMap := range d.GetLevelData(valStr) {
			if !isEqualData(validData, vMap.data) {
				compare := d.notesMap[vMap.fullPk]
				if compare == "" {
					compare = vMap.fullPk
				}
				return getMessageError(lang.Eq, d.message, validNotes, compare)
			}
		}
	} else {
		if !isEqualData(validData, val) {
			return getMessageError(lang.Eq, d.message, validNotes, val)
		}
	}
	return nil
}

func (m *methods) Gt(d *Data, args ...interface{}) error {
	return nil
}

func (m *methods) Gte(d *Data, args ...interface{}) error {
	return nil
}

func (m *methods) Lt(d *Data, args ...interface{}) error {
	return nil
}

func (m *methods) Lte(d *Data, args ...interface{}) error {
	return nil
}
