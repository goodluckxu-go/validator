package validator

// 数据
type Data struct {
	data          *interface{}
	notes         string
	fullField     string
	validData     interface{}
	ruleAsDataMap map[string]*ruleAsData
}

// 获取所有数据
func (d *Data) GetAllData() interface{} {
	return *d.data
}

// 获取数据且合并成数组
func (d *Data) GetData(key string) []interface{} {
	return getData(d.GetAllData(), key)
}

// 获取数组层级最近的一次相同数据(同一数组中)
func (d *Data) GetCommonData(key string) interface{} {
	commonField, _ := getCommonFullField(d.fullField, key)
	if d.ruleAsDataMap[commonField] == nil {
		return nil
	}
	return d.ruleAsDataMap[commonField].data
}

// 获取层级数据，遇到*合并数组
func (d *Data) GetLevelData(key string) []interface{} {
	commonField, otherField := getCommonFullField(d.fullField, key)
	if d.ruleAsDataMap[commonField] == nil {
		return nil
	}
	commonData := d.ruleAsDataMap[commonField].data
	list := getData(commonData, otherField)
	return list
}

// 获取验证数据
func (d *Data) GetValidData() interface{} {
	return d.validData
}

// 获取注释
func (d *Data) GetNotes() string {
	return d.notes
}

// 重设验证数据
func (d *Data) setValidData(value interface{}) {
	newData := *d.data
	ruleData := d.ruleAsDataMap[d.fullField]
	ruleData.data = value
	if d.fullField == "" {
		*d.data = value
		return
	}
	newData = upData(newData, d.fullField, value)
}
