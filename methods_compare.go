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
	if field, ok := val.(param.Field); ok {
		for _, vMap := range d.GetLevelData(string(field)) {
			if !isEqualData(validData, vMap.Data) {
				return validError(lang.Eq, d.getMessage(), langArg{
					notes:   d.GetNotes(),
					compare: d.GetNotesByPath(vMap.Path),
				})
			}
		}
	} else {
		if !isEqualData(validData, val) {
			return validError(lang.Eq, d.getMessage(), langArg{
				notes:   d.GetNotes(),
				compare: val,
			})
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
	if field, ok := val.(param.Field); ok {
		for _, vMap := range d.GetLevelData(string(field)) {
			compareData, err := isCompareData(validData, vMap.Data)
			if err != nil {
				return err
			}
			if compareData != 1 {
				return validError(lang.Gt, d.getMessage(), langArg{
					notes:   d.GetNotes(),
					compare: d.GetNotesByPath(vMap.Path),
				})
			}
		}
	} else {
		compareData, err := isCompareData(validData, val)
		if err != nil {
			return err
		}
		if compareData != 1 {
			return validError(lang.Gt, d.getMessage(), langArg{
				notes:   d.GetNotes(),
				compare: val,
			})
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
	if field, ok := val.(param.Field); ok {
		for _, vMap := range d.GetLevelData(string(field)) {
			compareData, err := isCompareData(validData, vMap.Data)
			if err != nil {
				return err
			}
			if compareData == -1 {
				return validError(lang.Gte, d.getMessage(), langArg{
					notes:   d.GetNotes(),
					compare: d.GetNotesByPath(vMap.Path),
				})
			}
		}
	} else {
		compareData, err := isCompareData(validData, val)
		if err != nil {
			return err
		}
		if compareData == -1 {
			return validError(lang.Gte, d.getMessage(), langArg{
				notes:   d.GetNotes(),
				compare: val,
			})
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
	if field, ok := val.(param.Field); ok {
		for _, vMap := range d.GetLevelData(string(field)) {
			compareData, err := isCompareData(validData, vMap.Data)
			if err != nil {
				return err
			}
			if compareData != -1 {
				return validError(lang.Lt, d.getMessage(), langArg{
					notes:   d.GetNotes(),
					compare: d.GetNotesByPath(vMap.Path),
				})
			}
		}
	} else {
		compareData, err := isCompareData(validData, val)
		if err != nil {
			return err
		}
		if compareData != -1 {
			return validError(lang.Lt, d.getMessage(), langArg{
				notes:   d.GetNotes(),
				compare: val,
			})
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
	if field, ok := val.(param.Field); ok {
		for _, vMap := range d.GetLevelData(string(field)) {
			compareData, err := isCompareData(validData, vMap.Data)
			if err != nil {
				return err
			}
			if compareData == 1 {
				return validError(lang.Lte, d.getMessage(), langArg{
					notes:   d.GetNotes(),
					compare: d.GetNotesByPath(vMap.Path),
				})
			}
		}
	} else {
		compareData, err := isCompareData(validData, val)
		if err != nil {
			return err
		}
		if compareData == 1 {
			return validError(lang.Lte, d.getMessage(), langArg{
				notes:   d.GetNotes(),
				compare: val,
			})
		}
	}
	return nil
}
