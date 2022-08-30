package condition

type Formula struct {
	Args []interface{}
}

// 括号
func Brackets(args ...interface{}) *Formula {
	return &Formula{Args: args}
}
