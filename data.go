package validator

import (
	"errors"
	"strings"
)

// Data 数据
type Data struct {
	data            *interface{}      // 原数据指针
	path            string            // 路径
	index           int               // 索引
	methodName      string            // 方法名称
	ruleTreeListPtr *[]ruleTree       // 数据指针
	messagesPtr     *[]Message        // 消息指针
	fileMapPtr      *map[string]*file // 文件指针
}

// GetAllData 获取所有数据
func (d *Data) GetAllData() interface{} {
	return *d.data
}

// GetData 获取数据且合并成数组
func (d *Data) GetData(path string) (rs []DataOne) {
	treeList := searchTree(path, *d.ruleTreeListPtr)
	rs = make([]DataOne, len(treeList))
	index := 0
	for _, v := range treeList {
		rs[index] = DataOne{
			Path: v.path,
			Data: v.data,
		}
		index++
	}
	return
}

// GetCommonData 获取数组层级最近的一次相同数据(同一数组中) todo
func (d *Data) GetCommonData(path string) interface{} {
	commonPath, otherPath := getCommonFullField(d.path, path)
	newPath := strings.TrimPrefix(commonPath+"."+otherPath, ".")
	treeList := searchTree(newPath, *d.ruleTreeListPtr)
	if len(treeList) > 0 {
		return treeList[0].data
	}
	return nil
}

// GetLevelData 获取层级数据，遇到*合并数组
func (d *Data) GetLevelData(path string) (rs []DataOne) {
	commonPath, otherPath := getCommonFullField(d.path, path)
	newPath := strings.TrimPrefix(commonPath+"."+otherPath, ".")
	return d.GetData(newPath)
}

// GetValidData 获取验证数据
func (d *Data) GetValidData() interface{} {
	return d.getValidData().data
}

// GetNotes 获取注释
func (d *Data) GetNotes() string {
	notes := d.getValidData().notes
	if notes == "" {
		notes = d.path
	}
	return notes
}

// GetNotesByPath 根据path获取注释
func (d *Data) GetNotesByPath(path string) string {
	rs := searchTree(path, *d.ruleTreeListPtr)
	if len(rs) > 0 {
		return rs[0].notes
	}
	return ""
}

// JumpValid 跳过当前字段验证
func (d *Data) JumpValid() error {
	return errors.New(jumpValid)
}

// JumpChild 跳过子集验证
func (d *Data) JumpChild() error {
	return errors.New(jumpChild)
}

func (d *Data) getFile() (f *file, err error) {
	if (*d.fileMapPtr)[d.path] == nil {
		return nil, errors.New("非文件类型")
	}
	return (*d.fileMapPtr)[d.path], nil
}

func (d *Data) getMessage() string {
	n := len(*d.messagesPtr)
	for i := n - 1; i >= 0; i-- {
		message := (*d.messagesPtr)[i]
		base := strings.TrimSuffix(message[0], stringJoin(".", "", d.methodName))
		if len(base) == len(message[0]) {
			continue
		}
		if inArrayRuleTree(d.path, searchTree(base, *d.ruleTreeListPtr)) {
			return message[1]
		}
	}
	return ""
}

func (d *Data) getValidData() ruleTree {
	return (*d.ruleTreeListPtr)[d.index]
}

func (d *Data) setValidData(value interface{}) {
	(*d.ruleTreeListPtr)[d.index].data = value
	*d.data = upData(*d.data, d.path, value)
}
