package validator

import "errors"

func (m *methods) Suffix(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, -1); err != nil {
		return err
	}
	if f, err := d.getFile(); err == nil {
		if !inArray(f.Suffix, args) {
			return validError(lang.Suffix, d.getMessage(), langArg{
				notes: d.GetNotes(),
				array: args,
			})
		}
	} else {
		return errors.New("只有文件可验证后缀")
	}
	return nil
}

func (m *methods) Mime(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, -1); err != nil {
		return err
	}
	if f, err := d.getFile(); err == nil {
		if !inArray(f.Mime, args) {
			return validError(lang.Mime, d.getMessage(), langArg{
				notes: d.GetNotes(),
				array: args,
			})
		}
	} else {
		return errors.New("只有文件可验证后缀")
	}
	return nil
}
