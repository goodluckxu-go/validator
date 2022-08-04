package types

import "test/validator/param"

var (
	Array   *param.Param // 数组类型
	Map     *param.Param // 对象类型
	String  *param.Param // 字符串类型
	Number  *param.Param // 数字类型
	Integer *param.Param // 整型
	Bool    *param.Param // 布尔型
	Date    *param.Param // 日期类型
	Field   *param.Param // 字段类型
)

func init() {
	Array = &param.Param{Type: 1}
	Map = &param.Param{Type: 2}
	String = &param.Param{Type: 3}
	Number = &param.Param{Type: 4}
	Integer = &param.Param{Type: 5}
	Bool = &param.Param{Type: 6}
	Date = &param.Param{Type: 7}
	Field = &param.Param{Type: 8}
}
