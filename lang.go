package validator

type language struct {
	Required string // 必填
	Array    string // 数组
	In       string // 数组中
	Email    string // 邮箱
	Phone    string // 手机
	Unique   string // 重复
	Regexp   string // 正则
	Map      string // 对象
	String   string // 字符串
	Number   string // 数字
	Integer  string // 整数
	Bool     string // 布尔
	Date     string // 日期
	Eq       string // 等于
	Gt       string // 大于
	Gte      string // 大于等于
	Lt       string // 小于
	Lte      string // 小于等于
	Len      string // 长度等于
	Min      string // 长度最低
	Max      string // 长度最高
}

func init() {
	lang = language{
		Required: "${notes}为必填",
		Array:    "${notes}必须是数组",
		Map:      "${notes}必须是对象",
		In:       "${notes}必须在数组${array}中",
		Email:    "${notes}必须是邮箱",
		Phone:    "${notes}必须是手机号",
		Unique:   "${notes}重复",
		Regexp:   "${notes}验证错误",
		String:   "${notes}必须是字符串",
		Number:   "${notes}必须是数字",
		Integer:  "${notes}必须是整数",
		Bool:     "${notes}必须是布尔",
		Date:     "${notes}必须是日期格式",
		Eq:       "${notes}必须等于${compare}",
		Gt:       "${notes}必须大于${compare}",
		Gte:      "${notes}必须大于等于${compare}",
		Lt:       "${notes}必须小于${compare}",
		Lte:      "${notes}必须小于等于${compare}",
		Len:      "${notes}长度必须是${len}",
		Min:      "${notes}最小长度为${len}",
		Max:      "${notes}最大长度为${len}",
	}
}
