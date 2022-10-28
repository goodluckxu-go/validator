package validator

import "errors"

func (m *methods) Suffix(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, -1); err != nil {
		return err
	}
	f := d.handle.fileMap[d.fullField]
	if f == nil {
		return errors.New("只有文件可验证后缀")
	}
	rsErr := getMessageError(lang.Suffix, d.message, d.GetNotes(), args)
	if !inArray(f.Suffix, args) {
		return rsErr
	}
	return nil
}

func (m *methods) Mime(d *Data, args ...interface{}) error {
	if err := validArgs(args, 1, -1); err != nil {
		return err
	}
	f := d.handle.fileMap[d.fullField]
	if f == nil {
		return errors.New("只有文件可验证类型")
	}
	rsErr := getMessageError(lang.Mime, d.message, d.GetNotes(), args)
	if !inArray(f.Mime, args) {
		return rsErr
	}
	return nil
}
