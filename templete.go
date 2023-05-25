package validator

import "net/http"

const (
	jumpValid string = "<**>.###(BREAK)###.<**>"    // 跳过该字段所有验证
	nextValid string = "<**>.###(CONTINUE)###.<**>" // 执行该字段下一个验证
)

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

// SetFun 设置自定义验证方法(其他参数可用于获取外部数据或者传地址修改外部数据)
func (m *method) SetFun(fn func(d *Data, args ...interface{}) error, args ...interface{}) (ms *method) {
	ms = getInstance(m).(*method)
	ms.methods = append(ms.methods, methodData{
		method: setMethodFunc(fn),
		args:   args,
	})
	return
}

// SetMethod 设置默认验证方法
func (m *method) SetMethod(r string, args ...interface{}) (ms *method) {
	ms = getInstance(m).(*method)
	ms.methods = append(ms.methods, methodData{
		method: r,
		args:   args,
	})
	return
}

// 存储
type storage struct {
	req      *http.Request // 请求
	data     interface{}   // 数据
	rules    []Rule        // 规则
	messages []Message     // 消息(覆盖规则使用)
}

// 处理数据
type handle struct {
	fileMap    map[string]*file
	ruleIndex  []string
	ruleData   map[string]*ruleAsData
	messages   [][3]string
	notesMap   map[string]string
	messageMap map[string]string
}

// 文件
type file struct {
	Suffix string // 后缀
	Mime   string // 协议
	Name   string // 文件名称
	Size   int64  // 文件大小
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
	data    interface{}
	methods []methodData
	notes   string
	isValid bool // 数据已经验证
}

// Rule 单个规则
type Rule struct {
	Field   string
	Methods []*method
	Notes   string
}

// Message 消息
type Message [2]string

// 单条数据
type DataOne struct {
	FullPk string
	Data   interface{}
}

// 是否过滤只获取验证参数
type filter bool
