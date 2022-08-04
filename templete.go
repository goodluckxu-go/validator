package validator

var (
	Method = method{
		nil, nil, 1, 2, 3, 4, 5, 6, 7, 8,
	} // 规则方法
	lang language // 语言
)

type fieldType uint8

// 规则方法
type method struct {
	methods []methodData
	list    []*method
	Array   fieldType // 数组类型
	Map     fieldType // 对象类型
	String  fieldType // 字符串类型
	Number  fieldType // 数字类型
	Integer fieldType // 整型
	Bool    fieldType // 布尔型
	Date    fieldType // 日期类型
	Field   fieldType // 字段类型
}

type methodData struct {
	method interface{}
	args   []interface{}
}

func (m *method) List(me ...*method) []*method {
	return me
}

// 设置自定义验证方法(其他参数可用于获取外部数据或者传地址修改外部数据)
func (m *method) SetFun(fn func(d *Data, args ...interface{}) error, args ...interface{}) (ms *method) {
	ms = getInstance(m).(*method)
	ms.methods = append(ms.methods, methodData{
		method: setMethodFunc(fn),
		args:   args,
	})
	return
}

// 设置默认验证方法
func (m *method) SetMethod(r string, args ...interface{}) (ms *method) {
	ms = getInstance(m).(*method)
	ms.methods = append(ms.methods, methodData{
		method: r,
		args:   args,
	})
	return
}

type value struct {
	value interface{}
}

// 设置规则类型
func (m *method) SetValue(v interface{}) value {
	return value{value: v}
}

// 处理数据
type ruleRow struct {
	field    string       // 验证字段(单个字段)
	pk       string       // 全字段
	methods  []methodData // 验证规则方法
	notes    string       // 字段注释
	children []ruleRow    // 子集
}

// 拆分数据组装
type ruleAsData struct {
	pk      string
	fullPk  string
	data    interface{}
	methods []methodData
	notes   string
}

// 单个规则
type Rule struct {
	Field   string
	Methods []*method
	Notes   string
}

// 消息
type Message [2]string

// 单条数据
type dataOne struct {
	fullPk string
	data   interface{}
}

// 是否过滤只获取验证参数
type filter bool
