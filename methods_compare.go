package validator

import (
	"github.com/goodluckxu-go/validator/param"
)

func (m *methods) Eq(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	val := args[0]
	validData := d.GetValidData()
	validNotes := d.GetNotes()
	if field, ok := val.(param.Field); ok {
		for _, vMap := range d.GetLevelData(string(field)) {
			if !isEqualData(validData, vMap.data) {
				compare := d.handle.notesMap[vMap.fullPk]
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
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	val := args[0]
	validData := d.GetValidData()
	validNotes := d.GetNotes()
	if field, ok := val.(param.Field); ok {
		for _, vMap := range d.GetLevelData(string(field)) {
			compareData, err := isCompareData(validData, vMap.data)
			if err != nil {
				return err
			}
			if compareData != 1 {
				compare := d.handle.notesMap[vMap.fullPk]
				if compare == "" {
					compare = vMap.fullPk
				}
				return getMessageError(lang.Gt, d.message, validNotes, compare)
			}
		}
	} else {
		compareData, err := isCompareData(validData, val)
		if err != nil {
			return err
		}
		if compareData != 1 {
			return getMessageError(lang.Gt, d.message, validNotes, val)
		}
	}
	return nil
}

func (m *methods) Gte(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	val := args[0]
	validData := d.GetValidData()
	validNotes := d.GetNotes()
	if field, ok := val.(param.Field); ok {
		for _, vMap := range d.GetLevelData(string(field)) {
			compareData, err := isCompareData(validData, vMap.data)
			if err != nil {
				return err
			}
			if compareData == -1 {
				compare := d.handle.notesMap[vMap.fullPk]
				if compare == "" {
					compare = vMap.fullPk
				}
				return getMessageError(lang.Gte, d.message, validNotes, compare)
			}
		}
	} else {
		compareData, err := isCompareData(validData, val)
		if err != nil {
			return err
		}
		if compareData == -1 {
			return getMessageError(lang.Gte, d.message, validNotes, val)
		}
	}
	return nil
}

func (m *methods) Lt(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	val := args[0]
	validData := d.GetValidData()
	validNotes := d.GetNotes()
	if field, ok := val.(param.Field); ok {
		for _, vMap := range d.GetLevelData(string(field)) {
			compareData, err := isCompareData(validData, vMap.data)
			if err != nil {
				return err
			}
			if compareData != -1 {
				compare := d.handle.notesMap[vMap.fullPk]
				if compare == "" {
					compare = vMap.fullPk
				}
				return getMessageError(lang.Lt, d.message, validNotes, compare)
			}
		}
	} else {
		compareData, err := isCompareData(validData, val)
		if err != nil {
			return err
		}
		if compareData != -1 {
			return getMessageError(lang.Lt, d.message, validNotes, val)
		}
	}
	return nil
}

func (m *methods) Lte(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, 1); err != nil {
		return err
	}
	val := args[0]
	validData := d.GetValidData()
	validNotes := d.GetNotes()
	if field, ok := val.(param.Field); ok {
		for _, vMap := range d.GetLevelData(string(field)) {
			compareData, err := isCompareData(validData, vMap.data)
			if err != nil {
				return err
			}
			if compareData == 1 {
				compare := d.handle.notesMap[vMap.fullPk]
				if compare == "" {
					compare = vMap.fullPk
				}
				return getMessageError(lang.Lte, d.message, validNotes, compare)
			}
		}
	} else {
		compareData, err := isCompareData(validData, val)
		if err != nil {
			return err
		}
		if compareData == 1 {
			return getMessageError(lang.Lte, d.message, validNotes, val)
		}
	}
	return nil
}
