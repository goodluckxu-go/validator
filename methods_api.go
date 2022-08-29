package validator

import (
	"sync"
)

var methodPool sync.Map

type methodFunc func(d *Data, args ...interface{}) error

type methodsApi interface {
	Required(d *Data, args ...interface{}) error
	ValidCondition(d *Data, args ...interface{}) error
	Array(d *Data, args ...interface{}) error
	Map(d *Data, args ...interface{}) error
	String(d *Data, args ...interface{}) error
	Number(d *Data, args ...interface{}) error
	Integer(d *Data, args ...interface{}) error
	Bool(d *Data, args ...interface{}) error
	Eq(d *Data, args ...interface{}) error
}

func newMethodsApi(api methodsApi) methodsApi {
	return api
}

func setMethodFunc(fn methodFunc) methodFunc {
	return fn
}
