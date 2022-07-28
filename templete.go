package validator

var (
	Method method   // 规则方法
	lang   language // 语言
)

// 规则方法
type method struct {
	methods []methodData
	list    []*method
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

// 注释
type Notes struct {
	notes string
}

// 是否过滤只获取验证参数
type filter bool
