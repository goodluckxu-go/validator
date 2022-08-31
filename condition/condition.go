package condition

type Formula struct {
	Args []interface{}
}

// Brackets 括号
func Brackets(args ...interface{}) *Formula {
	return &Formula{Args: args}
}
