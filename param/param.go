package param

// 传入类型
type Param struct {
	Type  uint8
	Value interface{}
}

func SetValue(v interface{}) *Param {
	return &Param{Value: v}
}
