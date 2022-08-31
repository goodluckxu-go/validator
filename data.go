package validator

// Data 数据
type Data struct {
	data          *interface{}
	notes         string
	fullField     string
	pk            string
	message       string
	validData     *interface{}
	notesMap      map[string]string      // 规则注释(notes)
	messageMap    map[string]string      // 规则注释(messages)
	ruleAsDataMap map[string]*ruleAsData // 数据
}

// GetAllData 获取所有数据
func (d *Data) GetAllData() interface{} {
	return *d.data
}

// GetData 获取数据且合并成数组
func (d *Data) GetData(key string) []dataOne {
	return getData(d.GetAllData(), key, "")
}

// GetCommonData 获取数组层级最近的一次相同数据(同一数组中)
func (d *Data) GetCommonData(key string) interface{} {
	commonField, _ := getCommonFullField(d.fullField, key)
	ruleData := d.ruleAsDataMap[commonField]
	if ruleData == nil {
		return nil
	}
	return ruleData.data
}

// GetLevelData 获取层级数据，遇到*合并数组
func (d *Data) GetLevelData(key string) []dataOne {
	commonField, otherField := getCommonFullField(d.fullField, key)
	ruleData := d.ruleAsDataMap[commonField]
	if ruleData == nil {
		return nil
	}
	list := getData(ruleData.data, otherField, "")
	return list
}

// GetValidData 获取验证数据
func (d *Data) GetValidData() interface{} {
	return *d.validData
}

// GetNotes 获取注释
func (d *Data) GetNotes() string {
	return d.notes
}

// 重设验证数据
func (d *Data) setValidData(value interface{}) {
	*d.validData = value
	newData := *d.data
	ruleData := d.ruleAsDataMap[d.fullField]
	ruleData.data = value
	if d.fullField == "" {
		*d.data = value
		return
	}
	newData = upData(newData, d.fullField, value)
	*d.data = newData
}
