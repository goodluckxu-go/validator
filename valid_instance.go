package validator

import (
	"net/http"
)

func init() {
	api := newMethodsApi(new(methods))
	methodPool.Store("valid_condition", setMethodFunc(api.ValidCondition))
	methodPool.Store("required", setMethodFunc(api.Required))
	methodPool.Store("nullable", setMethodFunc(api.Nullable))
	methodPool.Store("in", setMethodFunc(api.In))
	methodPool.Store("email", setMethodFunc(api.Email))
	methodPool.Store("phone", setMethodFunc(api.Phone))
	methodPool.Store("unique", setMethodFunc(api.Unique))
	methodPool.Store("regexp", setMethodFunc(api.Regexp))
	methodPool.Store("array", setMethodFunc(api.Array))
	methodPool.Store("map", setMethodFunc(api.Map))
	methodPool.Store("string", setMethodFunc(api.String))
	methodPool.Store("number", setMethodFunc(api.Number))
	methodPool.Store("integer", setMethodFunc(api.Integer))
	methodPool.Store("bool", setMethodFunc(api.Bool))
	methodPool.Store("date", setMethodFunc(api.Date))
	methodPool.Store("eq", setMethodFunc(api.Eq))
	methodPool.Store("gt", setMethodFunc(api.Gt))
	methodPool.Store("gte", setMethodFunc(api.Gte))
	methodPool.Store("lt", setMethodFunc(api.Lt))
	methodPool.Store("lte", setMethodFunc(api.Lte))
	methodPool.Store("len", setMethodFunc(api.Len))
	methodPool.Store("min", setMethodFunc(api.Min))
	methodPool.Store("max", setMethodFunc(api.Max))
}

// 实例化验证
func New(req *http.Request) *valid {
	v := new(valid)
	v.req = req
	return v
}

// 设置语言包地址
func SetLangAddr(langAddr string) {
	parseLang(langAddr)
}

// 设置全局验证规则方法
func RegisterMethod(key string, fn func(d *Data, args ...interface{}) error) {
	methodPool.Store(key, setMethodFunc(fn))
}
