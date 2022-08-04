package validator

import (
	"github.com/goodluckxu-go/validator/param"
	"sync"
)

var methodPool sync.Map

type methodFunc func(d *Data, args ...*param.Param) error

type methodsApi interface {
	Required(d *Data, args ...*param.Param) error
	Array(d *Data, args ...*param.Param) error
	Map(d *Data, args ...*param.Param) error
	String(d *Data, args ...*param.Param) error
	Number(d *Data, args ...*param.Param) error
	Integer(d *Data, args ...*param.Param) error
	Bool(d *Data, args ...*param.Param) error
	Eq(d *Data, args ...*param.Param) error

	ValidField(d *Data, args ...*param.Param) error
}

func newMethodsApi(api methodsApi) methodsApi {
	return api
}

func setMethodFunc(fn methodFunc) methodFunc {
	return fn
}
