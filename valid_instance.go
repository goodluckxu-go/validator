package validator

import (
	"github.com/goodluckxu-go/validator/param"
	"net/http"
)

func init() {
	api := newMethodsApi(new(methods))
	methodPool.Store("required", setMethodFunc(api.Required))
	methodPool.Store("array", setMethodFunc(api.Array))
	methodPool.Store("map", setMethodFunc(api.Map))
	methodPool.Store("string", setMethodFunc(api.String))
	methodPool.Store("number", setMethodFunc(api.Number))
	methodPool.Store("integer", setMethodFunc(api.Integer))
	methodPool.Store("bool", setMethodFunc(api.Bool))
	methodPool.Store("eq", setMethodFunc(api.Eq))
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
	methodPool.Store(key, setMethodFunc(func(d *Data, args ...*param.Param) error {
		var newArgs []interface{}
		for _, arg := range args {
			newArgs = append(newArgs, arg.Value)
		}
		return fn(d, newArgs...)
	}))
}
