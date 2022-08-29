package validator

import (
	"fmt"
	"test/validator/param"
)

func (m *methods) Eq(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	val := args[0]
	validData := d.GetValidData()
	validNotes := d.GetNotes()
	if file, ok := val.(param.File); ok {
		for _, vMap := range d.GetLevelData(string(file)) {
			fmt.Println(validData, vMap.data)
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
