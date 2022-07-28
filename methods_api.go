package validator

import "sync"

var methodPool sync.Map

type methodFunc func(d *Data, args ...interface{}) error

type methodsApi interface {
	Required(d *Data, args ...interface{}) error
	ValidField(d *Data, args ...interface{}) error
}

func newMethodsApi(api methodsApi) methodsApi {
	return api
}

func setMethodFunc(fn methodFunc) methodFunc {
	return fn
}
