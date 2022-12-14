package validator

func init() {
	api := newMethodsApi(new(methods))
	methodPool.Store("valid_condition", setMethodFunc(api.ValidCondition))
	methodPool.Store("required", setMethodFunc(api.Required))
	methodPool.Store("nullable", setMethodFunc(api.Nullable))
	methodPool.Store("in", setMethodFunc(api.In))
	methodPool.Store("not_in", setMethodFunc(api.NotIn))
	methodPool.Store("email", setMethodFunc(api.Email))
	methodPool.Store("phone", setMethodFunc(api.Phone))
	methodPool.Store("unique", setMethodFunc(api.Unique))
	methodPool.Store("regexp", setMethodFunc(api.Regexp))
	methodPool.Store("not_regexp", setMethodFunc(api.NotRegexp))
	methodPool.Store("array", setMethodFunc(api.Array))
	methodPool.Store("map", setMethodFunc(api.Map))
	methodPool.Store("string", setMethodFunc(api.String))
	methodPool.Store("number", setMethodFunc(api.Number))
	methodPool.Store("integer", setMethodFunc(api.Integer))
	methodPool.Store("bool", setMethodFunc(api.Bool))
	methodPool.Store("date", setMethodFunc(api.Date))
	methodPool.Store("file", setMethodFunc(api.File))
	methodPool.Store("eq", setMethodFunc(api.Eq))
	methodPool.Store("gt", setMethodFunc(api.Gt))
	methodPool.Store("gte", setMethodFunc(api.Gte))
	methodPool.Store("lt", setMethodFunc(api.Lt))
	methodPool.Store("lte", setMethodFunc(api.Lte))
	methodPool.Store("len", setMethodFunc(api.Len))
	methodPool.Store("min", setMethodFunc(api.Min))
	methodPool.Store("max", setMethodFunc(api.Max))
	methodPool.Store("suffix", setMethodFunc(api.Suffix))
	methodPool.Store("mime", setMethodFunc(api.Mime))
}

// New 实例化验证
func New() *Valid {
	v := new(Valid)
	v.handle.fileMap = map[string]*file{}
	v.handle.ruleData = map[string]*ruleAsData{}
	v.handle.notesMap = map[string]string{}
	v.handle.messageMap = map[string]string{}
	return v
}

// SetLangAddr 设置语言包地址
func SetLangAddr(langAddr string) {
	parseLang(langAddr)
}

// RegisterMethod 设置全局验证规则方法
func RegisterMethod(key string, fn func(d *Data, args ...interface{}) error) {
	methodPool.Store(key, setMethodFunc(fn))
}
