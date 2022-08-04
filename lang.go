package validator

type language struct {
	Required string // 必填
	Array    string
	Map      string
	String   string
	Number   string
	Integer  string
	Bool     string
	Eq       string
}

func init() {
	lang = language{
		Required: "${notes}为必填",
		Array:    "${notes}必须是数组",
		Map:      "${notes}必须是对象",
		String:   "${notes}必须是字符串",
		Number:   "${notes}必须是数字",
		Integer:  "${notes}必须是整数",
		Bool:     "${notes}必须是布尔",
		Eq:       "${notes}必须等于${compare}",
	}
}
