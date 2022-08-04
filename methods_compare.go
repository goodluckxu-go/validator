package validator

func (m *methods) Eq(d *Data, args ...interface{}) error {
	var fType fieldType
	var val value
	for _, arg := range args {
		switch arg.(type) {
		case fieldType:
			fType = arg.(fieldType)
		case value:
			val = arg.(value)
		}
	}
	validData := d.GetValidData()
	validNotes := d.GetNotes()
	if fType == Method.Field {
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
