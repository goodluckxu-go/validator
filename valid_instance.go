package validator

func init() {
	api := methodsApi(new(methods))
	methodPool.Store("valid_condition", methodFunc(api.ValidCondition))
	methodPool.Store("required", methodFunc(api.Required))
	methodPool.Store("nullable", methodFunc(api.Nullable))
	methodPool.Store("in", methodFunc(api.In))
	methodPool.Store("not_in", methodFunc(api.NotIn))
	methodPool.Store("email", methodFunc(api.Email))
	methodPool.Store("phone", methodFunc(api.Phone))
	methodPool.Store("unique", methodFunc(api.Unique))
	methodPool.Store("regexp", methodFunc(api.Regexp))
	methodPool.Store("not_regexp", methodFunc(api.NotRegexp))
	methodPool.Store("array", methodFunc(api.Array))
	methodPool.Store("map", methodFunc(api.Map))
	methodPool.Store("string", methodFunc(api.String))
	methodPool.Store("number", methodFunc(api.Number))
	methodPool.Store("integer", methodFunc(api.Integer))
	methodPool.Store("bool", methodFunc(api.Bool))
	methodPool.Store("date", methodFunc(api.Date))
	methodPool.Store("file", methodFunc(api.File))
	methodPool.Store("eq", methodFunc(api.Eq))
	methodPool.Store("gt", methodFunc(api.Gt))
	methodPool.Store("gte", methodFunc(api.Gte))
	methodPool.Store("lt", methodFunc(api.Lt))
	methodPool.Store("lte", methodFunc(api.Lte))
	methodPool.Store("len", methodFunc(api.Len))
	methodPool.Store("min", methodFunc(api.Min))
	methodPool.Store("max", methodFunc(api.Max))
	methodPool.Store("suffix", methodFunc(api.Suffix))
	methodPool.Store("mime", methodFunc(api.Mime))
}

// New 实例化验证
func New() *Valid {
	v := new(Valid)
	v.handle.fileMap = map[string]*file{}
	return v
}

// SetLangAddr 设置语言包地址
func SetLangAddr(langAddr string) {
	parseLang(langAddr)
}

// RegisterMethod 设置全局验证规则方法
func RegisterMethod(key string, fn func(d *Data, args ...interface{}) error) {
	methodPool.Store(key, methodFunc(fn))
}
