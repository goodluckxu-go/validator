package validator

import (
	"sync"
)

var methodPool sync.Map

type methodFunc func(d *Data, args ...interface{}) error

type methodsApi interface {
	// 通用验证
	ValidCondition(d *Data, args ...interface{}) error
	Required(d *Data, args ...interface{}) error
	Nullable(d *Data, args ...interface{}) error
	In(d *Data, args ...interface{}) error
	NotIn(d *Data, args ...interface{}) error
	Unique(d *Data, args ...interface{}) error
	Email(d *Data, args ...interface{}) error
	Phone(d *Data, args ...interface{}) error
	Regexp(d *Data, args ...interface{}) error
	NotRegexp(d *Data, args ...interface{}) error
	// 类型验证
	Array(d *Data, args ...interface{}) error
	Map(d *Data, args ...interface{}) error
	String(d *Data, args ...interface{}) error
	Number(d *Data, args ...interface{}) error
	Integer(d *Data, args ...interface{}) error
	Bool(d *Data, args ...interface{}) error
	Date(d *Data, args ...interface{}) error
	File(d *Data, args ...interface{}) error
	// 比较验证
	Eq(d *Data, args ...interface{}) error
	Gt(d *Data, args ...interface{}) error
	Gte(d *Data, args ...interface{}) error
	Lt(d *Data, args ...interface{}) error
	Lte(d *Data, args ...interface{}) error
	// 长度验证
	Len(d *Data, args ...interface{}) error
	Min(d *Data, args ...interface{}) error
	Max(d *Data, args ...interface{}) error
	// 文件验证
	Suffix(d *Data, args ...interface{}) error
	Mime(d *Data, args ...interface{}) error
}
