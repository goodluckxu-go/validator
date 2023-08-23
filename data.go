package validator

import (
	"errors"
	"strings"
)

// Data 数据
type Data struct {
	data           *interface{}      // 原数据指针
	path           string            // 路径
	methodName     string            // 方法名称
	ruleRowListPtr *[]ruleRow        // 数据指针
	pathIndexPtr   *map[string]int   // 验证数据索引
	messagesPtr    *[]Message        // 消息指针
	fileMapPtr     *map[string]*file // 文件指针
}

// GetAllData 获取所有数据
func (d *Data) GetAllData() interface{} {
	return *d.data
}

// GetData 获取数据且合并成数组
func (d *Data) GetData(path string) (rs []DataOne) {
	for _, v := range *d.ruleRowListPtr {
		if inArrayString(path, v.samePaths) {
			rs = append(rs, DataOne{
				FullPk: v.path,
				Data:   v.data,
			})
		}
	}
	return
}

// GetCommonData 获取数组层级最近的一次相同数据(同一数组中) todo
func (d *Data) GetCommonData(path string) interface{} {
	commonPath, otherPath := getCommonFullField(d.path, path)
	newPath := strings.TrimPrefix(commonPath+"."+otherPath, ".")
	pathIndex := *d.pathIndexPtr
	if index, ok := pathIndex[newPath]; ok {
		ruleRowList := *d.ruleRowListPtr
		return ruleRowList[index].data
	}
	return nil
}

// GetLevelData 获取层级数据，遇到*合并数组
func (d *Data) GetLevelData(path string) (rs []DataOne) {
	commonPath, _ := getCommonFullField(d.path, path)
	for _, v := range *d.ruleRowListPtr {
		if strings.Index(v.path, commonPath) != -1 && inArrayString(path, v.samePaths) {
			rs = append(rs, DataOne{
				FullPk: v.path,
				Data:   v.data,
			})
		}
	}
	return
}

// GetValidData 获取验证数据
func (d *Data) GetValidData() interface{} {
	return d.getValidData().data
}

// SetValidData 重设验证数据
func (d *Data) SetValidData(value interface{}) {
	(*d.ruleRowListPtr)[(*d.pathIndexPtr)[d.path]].data = value
	*d.data = upData(*d.data, d.path, value)
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
	if index, ok := (*d.pathIndexPtr)[path]; ok {
		notes := (*d.ruleRowListPtr)[index].notes
		if notes == "" {
			notes = (*d.ruleRowListPtr)[index].path
		}
		return notes
	}
	return ""
}

// JumpValid 跳过当前字段验证
func (d *Data) JumpValid() error {
	return errors.New(jumpValid)
}

func (d *Data) getFile() (f *file, err error) {
	if (*d.fileMapPtr)[d.path] == nil {
		return nil, errors.New("非文件类型")
	}
	return (*d.fileMapPtr)[d.path], nil
}

func (d *Data) getMessage() string {
	for i := len(*d.messagesPtr) - 1; i >= 0; i-- {
		message := (*d.messagesPtr)[i]
		base := strings.TrimSuffix(message[0], stringJoin(".", "", d.methodName))
		if len(base) == len(message[0]) {
			continue
		}
		if inArrayString(base, d.getValidData().samePaths) {
			return message[1]
		}
	}
	return ""
}

func (d *Data) getValidData() ruleRow {
	return (*d.ruleRowListPtr)[(*d.pathIndexPtr)[d.path]]
}
