package validator

import (
	"net/http"
)

const (
	jumpValid string = "<**>.###(BREAK)###.<**>" // 跳过该字段所有验证
	jumpChild string = "<**>.###(CHILD)###.<**>" // 跳过子集验证
	noPrefix  string = "no prefix"               // 不存在前缀
)

var (
	Method method   // 规则方法
	lang   language // 语言
)

// 规则方法
type method struct {
	methods []methodData
}

type methodData struct {
	method interface{}
	args   []interface{}
}

func (m *method) List(me ...*method) []methodData {
	size := 0
	for _, mm := range me {
		size += len(mm.methods)
	}
	rs := make([]methodData, size)
	index := 0
	for _, mm := range me {
		for _, v := range mm.methods {
			rs[index] = v
			index++
		}
	}
	return rs
}

// SetFun 设置自定义验证方法(其他参数可用于获取外部数据或者传地址修改外部数据)
func (m *method) SetFun(fn func(d *Data, args ...interface{}) error, args ...interface{}) (ms *method) {
	ms = m.getInstance()
	ms.methods = append(ms.methods, methodData{
		method: methodFunc(fn),
		args:   args,
	})
	return
}

// SetMethod 设置默认验证方法
func (m *method) SetMethod(r string, args ...interface{}) (ms *method) {
	ms = m.getInstance()
	ms.methods = append(ms.methods, methodData{
		method: r,
		args:   args,
	})
	return
}

func (m *method) getInstance() (ms *method) {
	ms = new(method)
	ms.methods = m.methods
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
	fileMap          map[string]*file // 文件数据
	ruleTreeList     []ruleTree       // 规则树
	validRulePtrList []*ruleTree      // 验证规则列表
}

// 文件
type file struct {
	Suffix string // 后缀
	Mime   string // 协议
	Name   string // 文件名称
	Size   int64  // 文件大小
}

type langArg struct {
	notes   interface{} // 注释
	array   interface{} // 数组
	compare interface{} // 比较
	len     interface{} // 长度
}

// 规则和数据对应的单条规则树
type ruleTree struct {
	path            string       // 路径
	prefix          string       // 前缀
	field           string       // 当前字段
	index           int          // 索引
	parentIndex     int          // 父级索引
	firstChildIndex int          // 第一个孩子索引
	lastChildIndex  int          // 最后一个孩子索引
	data            interface{}  // 数据
	notes           string       // 字段注释
	methods         []methodData // 验证规则方法
	jumpChild       bool         // 是否跳过子集验证
}

// Rule 单个规则
type Rule struct {
	prefix       string      // 临时前缀
	data         interface{} // 临时数据
	samePrefixes []string    // 临时相似前缀
	Field        string
	Methods      []methodData
	Notes        string
}

// Message 消息
type Message [2]string

// DataOne 单条数据
type DataOne struct {
	Path string
	Data interface{}
}

// 是否过滤只获取验证参数
type filter bool
