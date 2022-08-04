package validator

import (
	"github.com/goodluckxu-go/validator/param"
	"github.com/goodluckxu-go/validator/types"
)

func (m *methods) Eq(d *Data, args ...*param.Param) error {
	var fType *param.Param
	var val interface{}
	for _, arg := range args {
		switch arg {
		case types.Field:
			fType = arg
		default:
			if arg.Value != nil {
				val = arg.Value
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

func (m *methods) Gt(d *Data, args ...param.Param) error {
	return nil
}

func (m *methods) Gte(d *Data, args ...param.Param) error {
	return nil
}

func (m *methods) Lt(d *Data, args ...param.Param) error {
	return nil
}

func (m *methods) Lte(d *Data, args ...param.Param) error {
	return nil
}
