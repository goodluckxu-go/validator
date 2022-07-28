package validator

import (
	"net/http"
)

func init() {
	api := newMethodsApi(new(methods))
	methodPool.Store("required", setMethodFunc(api.Required))
	methodPool.Store("valid_field", setMethodFunc(api.ValidField))
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
