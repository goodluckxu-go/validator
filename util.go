package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goodluckxu-go/validator/condition"
	"github.com/goodluckxu-go/validator/param"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 读取body内容
func readBody(r *http.Request) []byte {
	// 从原有Request.Body读取
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return bodyBytes
	}
	// 新建缓冲区并替换原有Request.body
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

// 将规则和数据处理成单条
func disintegrateRules(rules []Rule, data interface{}, init bool, rsPtr *[]ruleTree, args ...interface{}) (err error) {
	rs := *rsPtr
	pathIndexMap := map[string]int{}
	pathIndex := 0
	if len(args) > 1 {
		pathIndexMap, _ = args[0].(map[string]int)
		pathIndex, _ = args[1].(int)
	}
	var otherRules []Rule
	boolMap := map[string]interface{}{} // 是否时map类型
	n := len(rules)
	sameBeforePrefix := ""
	for index := 0; index < n; index++ {
		rule := rules[index]
		if !init {
			data = rule.data
		}
		path := ""
		if _, ok := pathIndexMap[path]; !ok {
			pathIndexMap[path] = pathIndex
			pathIndex++
			rs = append(rs, ruleTree{
				path:        path,
				prefix:      noPrefix,
				data:        data,
				index:       pathIndex - 1,
				parentIndex: -1,
			})
		}
		if rule.prefix == "" && rule.Field == "" {
			rs[0].methods = rule.Methods
			rs[0].notes = rule.Notes
		}
		fieldList := strings.Split(rule.Field, ".")
		if fieldList[0] == "*" {
			if boolMap[rule.prefix] == nil {
				boolMap[rule.prefix] = false
			} else if boolMap[rule.prefix] == true {
				err = fmt.Errorf("%s冲突，map和slice不能并存", rule.prefix)
				return
			}
			// 数组
			if sameBeforePrefix == rule.prefix {
				continue
			}
			sameBeforePrefix = rule.prefix
			jump := 0
			for i := index; i < n; i++ {
				newRule := rules[i]
				if newRule.prefix != sameBeforePrefix {
					break
				}
				jump++
			}
			dataList, _ := data.([]interface{})
			// 带*的变成0
			if len(dataList) == 0 {
				dataList = []interface{}{nil}
			}
			for k, v := range dataList {
				for i := index; i < index+jump; i++ {
					newRule := rules[i]
					newFieldList := strings.Split(newRule.Field, ".")
					if len(newFieldList) == 1 {
						continue
					}
					if newFieldList[1] == "*" {
						otherRules = append(otherRules, Rule{
							Field:   strings.Join(newFieldList[1:], "."),
							data:    v,
							prefix:  fieldJoin(rule.prefix, k),
							Methods: newRule.Methods,
							Notes:   newRule.Notes,
						})
					} else {
						vMap, _ := v.(map[string]interface{})
						otherRules = append(otherRules, Rule{
							Field: strings.Join(newFieldList[1:], "."),
							data: map[string]interface{}{
								newFieldList[1]: vMap[newFieldList[1]],
							},
							prefix:  fieldJoin(rule.prefix, k),
							Methods: newRule.Methods,
							Notes:   newRule.Notes,
						})
					}
				}
				path = fieldJoin(rule.prefix, k)
				pathIndexMap[path] = pathIndex
				pathIndex++
				parentIndex := -1
				if tmp, ok := pathIndexMap[rule.prefix]; ok {
					parentIndex = tmp
					rsLen := len(rs)
					if rs[parentIndex].firstChildIndex == 0 {
						rs[parentIndex].firstChildIndex = rsLen
					}
					rs[parentIndex].lastChildIndex = rsLen
				}
				treeItem := ruleTree{
					path:        path,
					prefix:      rule.prefix,
					field:       fieldJoin(k),
					data:        v,
					index:       pathIndex - 1,
					parentIndex: parentIndex,
				}
				if len(fieldList) == 1 {
					treeItem.methods = rule.Methods
					treeItem.notes = rule.Notes
				}
				rs = append(rs, treeItem)
			}
			index += jump - 1
			continue
		} else {
			if boolMap[rule.prefix] == nil {
				boolMap[rule.prefix] = true
			} else if boolMap[rule.prefix] == false {
				err = fmt.Errorf("%s冲突，map和slice不能并存", rule.prefix)
				return
			}
			// 对象
			dataMap, _ := data.(map[string]interface{})
			if len(fieldList) > 1 {
				otherRules = append(otherRules, Rule{
					Field:   strings.Join(fieldList[1:], "."),
					data:    dataMap[fieldList[0]],
					prefix:  fieldJoin(rule.prefix, fieldList[0]),
					Methods: rule.Methods,
					Notes:   rule.Notes,
				})
			}
			path = fieldJoin(rule.prefix, fieldList[0])
			if _, ok := pathIndexMap[path]; ok {
				continue
			}
			pathIndexMap[path] = pathIndex
			pathIndex++
			parentIndex := -1
			if tmp, ok := pathIndexMap[rule.prefix]; ok {
				parentIndex = tmp
				rsLen := len(rs)
				if rs[parentIndex].firstChildIndex == 0 {
					rs[parentIndex].firstChildIndex = rsLen
				}
				rs[parentIndex].lastChildIndex = rsLen
			}
			treeItem := ruleTree{
				path:        path,
				prefix:      rule.prefix,
				field:       fieldList[0],
				data:        dataMap[fieldList[0]],
				index:       pathIndex - 1,
				parentIndex: parentIndex,
			}
			if len(fieldList) == 1 {
				treeItem.methods = rule.Methods
				treeItem.notes = rule.Notes
			}
			rs = append(rs, treeItem)
		}
	}
	if len(otherRules) > 0 {
		if err = disintegrateRules(otherRules, nil, false, &rs, pathIndexMap, pathIndex); err != nil {
			return
		}
	}
	*rsPtr = rs
	return
}

func ruleTreeSort(list []ruleTree, index int) (rs []*ruleTree) {
	size := len(list)
	stackList := make([]*ruleTree, size)
	stackIndex := 0
	rs = make([]*ruleTree, size)
	rsIndex := 0
	n := len(list)
	for i := 0; i < n; i++ {
		if list[i].parentIndex == index {
			stackList[stackIndex] = &list[i]
			stackIndex++
			for stackIndex > 0 {
				stackIndex--
				curr := stackList[stackIndex]
				rs[rsIndex] = curr
				rsIndex++
				if curr.lastChildIndex == 0 {
					continue
				}
				for j := curr.lastChildIndex; j >= curr.firstChildIndex; j-- {
					stackList[stackIndex] = &list[j]
					stackIndex++
				}
			}
		} else {
			break
		}
	}
	return
}

type stack struct {
	path string
	tree ruleTree
}

func searchTree(path string, treeList []ruleTree) (rs []ruleTree) {
	var stackList []stack
	stackList = append(stackList, stack{
		path: path,
		tree: treeList[0],
	})
	for len(stackList) > 0 {
		curr := stackList[len(stackList)-1]
		stackList = stackList[:len(stackList)-1]
		// 搜索
		pathList := strings.Split(curr.path, ".")
		for i := curr.tree.firstChildIndex; i <= curr.tree.lastChildIndex; i++ {
			newTree := treeList[i]
			if pathList[0] == "*" {
				if len(pathList) == 1 {
					rs = append(rs, newTree)
				} else {
					stackList = append(stackList, stack{
						path: strings.Join(pathList[1:], "."),
						tree: newTree,
					})
				}
			} else if treeList[i].field == pathList[0] {
				if len(pathList) == 1 {
					rs = append(rs, newTree)
				} else {
					stackList = append(stackList, stack{
						path: strings.Join(pathList[1:], "."),
						tree: newTree,
					})
				}
				break
			}
		}
	}
	// 倒序
	n := len(rs)
	for i := 0; i < n/2; i++ {
		rs[i], rs[n-1-i] = rs[n-1-i], rs[i]
	}
	return rs
}

func isSlice(prefix string) bool {
	if prefix == "*" {
		return true
	}
	if _, err := strconv.Atoi(prefix); err == nil {
		return true
	}
	return false
}

func fieldJoin(vs ...interface{}) string {
	if len(vs) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for _, v := range vs {
		s := ""
		switch v.(type) {
		case string:
			s = v.(string)
		case int:
			s = strconv.Itoa(v.(int))
		}
		if s == "" {
			continue
		}
		if buf.Len() > 0 {
			_, _ = buf.WriteString(".")
		}
		_, _ = buf.WriteString(s)
	}
	return buf.String()
}

func stringJoin(s, deletePrefix string, vs ...interface{}) string {
	var buf bytes.Buffer
	if s != "" {
		buf.WriteString(s)
	}
	for _, v := range vs {
		switch v.(type) {
		case string:
			buf.WriteString(v.(string))
		case []byte:
			buf.Write(v.([]byte))
		default:
			buf.WriteString(toString(v))
		}
	}
	s = buf.String()
	if deletePrefix != "" {
		s = strings.TrimPrefix(s, deletePrefix)
	}
	return s
}

func upData(data interface{}, key string, value interface{}) interface{} {
	if key == "" {
		return value
	}
	keyList := strings.Split(key, ".")
	firstKey := keyList[0]
	otherKey := strings.Join(keyList[1:], ".")
	switch data.(type) {
	case map[string]interface{}:
		dataMap, _ := data.(map[string]interface{})
		dataMap[firstKey] = upData(dataMap[firstKey], otherKey, value)
		data = dataMap
	case []interface{}:
		index, _ := strconv.Atoi(firstKey)
		dataList, _ := data.([]interface{})
		if len(dataList) > index {
			dataList[index] = upData(dataList[index], otherKey, value)
		}
		data = dataList
	}
	return data
}

func parseLang(langAddr string) {
	if langAddr == "" {
		return
	}
	var f *os.File
	var err error
	if f, err = os.Open(langAddr); err != nil {
		panic(err)
	}
	var by []byte
	if by, err = io.ReadAll(f); err != nil {
		panic(err)
	}
	langMap := map[string]interface{}{}
	if err = json.Unmarshal(by, &langMap); err != nil {
		panic(err)
	}
	langValue := reflect.ValueOf(&lang).Elem()
	langLen := langValue.NumField()
	for key, val := range langMap {
		for i := 0; i < langLen; i++ {
			field := langValue.Type().Field(i).Name
			if key == field {
				langValue.Field(i).Set(reflect.ValueOf(val))
				break
			}
		}
	}
}

func validError(langStr string, message string, arg langArg) error {
	if message != "" {
		return errors.New(message)
	}
	if arg.notes != nil {
		langStr = strings.ReplaceAll(langStr, "${notes}", toString(arg.notes))
	}
	if arg.array != nil {
		langStr = strings.ReplaceAll(langStr, "${array}", toString(arg.array))
	}
	if arg.compare != nil {
		langStr = strings.ReplaceAll(langStr, "${compare}", toString(arg.compare))
	}
	if arg.len != nil {
		langStr = strings.ReplaceAll(langStr, "${len}", toString(arg.len))
	}
	return errors.New(langStr)
}

func getCommonFullField(field, otherField string) (string, string) {
	fieldList := strings.Split(field, ".")
	otherFieldList := strings.Split(otherField, ".")
	var rsList []string
	i := 0
	for index, ov := range otherFieldList {
		if len(fieldList) <= index {
			break
		}
		fv := fieldList[index]
		if ov != "*" && ov != fv {
			break
		}
		rsList = append(rsList, fv)
		i++
	}
	return strings.Join(rsList, "."), strings.Join(otherFieldList[i:], ".")
}

// 比较两个数是否相等
func isEqualData(dataOne, dataTwo interface{}) bool {
	switch dataOne.(type) {
	case string:
		if dataTwoDate, ok := dataTwo.(time.Time); ok {
			dataOneDate, err := timeParse(fmt.Sprintf("%v", dataOne))
			if err != nil {
				return false
			}
			if dataOneDate.Unix() == dataTwoDate.Unix() {
				return true
			}
		}
		if dataOne.(string) == fmt.Sprintf("%v", dataTwo) {
			return true
		}
	case int, float64:
		dataOneFloat64 := toFloat64(dataOne)
		dataTwoFloat64 := toFloat64(dataTwo)
		if dataOneFloat64 == dataTwoFloat64 {
			return true
		}
	case bool:
		if _, ok := dataTwo.(bool); !ok {
			return false
		}
		if compareDataBool, bl := dataTwo.(bool); bl {
			if dataOne.(bool) == compareDataBool {
				return true
			}
		}
	default:
		if reflect.DeepEqual(dataOne, dataTwo) {
			return true
		}
	}
	return false
}

// 比较两个数值，数字或字符串的数字或日期
func isCompareData(dataOne, dataTwo interface{}) (int, error) {
	var dataOneFloat64, dataTwoFloat64 float64
	var err error
	errInfo := errors.New("比较的数必须是数字,字符串的数字或日期且类型匹配")
	switch dataOne.(type) {
	case int64, float64:
		dataOneFloat64 = toFloat64(dataOne)
		dataTwoFloat64 = toFloat64(dataTwo)
	case string:
		isNumber := true
		dataOneFloat64, err = strconv.ParseFloat(fmt.Sprintf("%v", dataOne), 64)
		if err != nil {
			isNumber = false
		}
		if isNumber {
			dataTwoFloat64, err = strconv.ParseFloat(fmt.Sprintf("%v", dataTwo), 64)
			if err != nil {
				return 0, errInfo
			}
		} else {
			var dataOneDate, dataTwoDate time.Time
			dataOneDate, err = timeParse(fmt.Sprintf("%v", dataOne))
			if err != nil {
				return 0, errInfo
			}
			var ok bool
			if dataTwoDate, ok = dataTwo.(time.Time); !ok {
				dataTwoDate, err = timeParse(fmt.Sprintf("%v", dataTwo))
				if err != nil {
					return 0, errInfo
				}
			}
			dataOneFloat64 = float64(dataOneDate.Unix())
			dataTwoFloat64 = float64(dataTwoDate.Unix())
		}
	default:
		return 0, errInfo
	}
	if dataOneFloat64 > dataTwoFloat64 {
		return 1, nil
	}
	if dataOneFloat64 < dataTwoFloat64 {
		return -1, nil
	}
	return 0, nil
}

// 字符串转时间
func timeParse(date string) (time.Time, error) {
	formatAtByte := []byte("0000-00-00 00:00:00")
	copy(formatAtByte, date)
	return time.ParseInLocation("2006-01-02 15:04:05", string(formatAtByte), time.Local)
}

/**
 * 验证时间
 */
func validDate(date string, format string) (err error) {
	err = fmt.Errorf("日期格式和值不匹配，格式：%s, 值：%s", format, date)
	format = strings.ReplaceAll(format, "Y", "YYYY")
	format = strings.ReplaceAll(format, "m", "mm")
	format = strings.ReplaceAll(format, "d", "dd")
	format = strings.ReplaceAll(format, "H", "HH")
	format = strings.ReplaceAll(format, "i", "ii")
	format = strings.ReplaceAll(format, "s", "ss")
	if len(date) != len(format) {
		return
	}
	if err = validSingleDate("YYYY", &date, &format); err != nil {
		return err
	}
	if err = validSingleDate("mm", &date, &format); err != nil {
		return err
	}
	if err = validSingleDate("dd", &date, &format); err != nil {
		return err
	}
	if err = validSingleDate("HH", &date, &format); err != nil {
		return err
	}
	if err = validSingleDate("ii", &date, &format); err != nil {
		return err
	}
	if err = validSingleDate("ss", &date, &format); err != nil {
		return err
	}
	if date != format {
		return
	}
	return nil
}

func validSingleDate(single string, date, format *string) error {
	lenSingle := len(single)
	newDate := *date
	newFormat := *format
	for {
		index := strings.Index(newFormat, single)
		if index == -1 {
			break
		}
		validDate := newDate[index : index+lenSingle]
		switch single {
		case "YYYY":
			if !regexp.MustCompile(`\d{4}`).MatchString(validDate) {
				return fmt.Errorf("年格式和值不匹配，格式：Y, 值：%s", validDate)
			}
		case "mm":
			validInt64, _ := strconv.ParseInt(validDate, 10, 64)
			if validInt64 < 1 || validInt64 > 12 {
				return fmt.Errorf("月格式和值不匹配，格式：m, 值：%s", validDate)
			}
		case "dd":
			validInt64, _ := strconv.ParseInt(validDate, 10, 64)
			if validInt64 < 1 || validInt64 > 31 {
				return fmt.Errorf("日格式和值不匹配，格式：d, 值：%s", validDate)
			}
		case "HH":
			validInt64, _ := strconv.ParseInt(validDate, 10, 64)
			if validInt64 < 0 || validInt64 > 23 {
				return fmt.Errorf("时格式和值不匹配，格式：H, 值：%s", validDate)
			}
		case "ii":
			validInt64, _ := strconv.ParseInt(validDate, 10, 64)
			if validInt64 < 0 || validInt64 > 59 {
				return fmt.Errorf("分格式和值不匹配，格式：i, 值：%s", validDate)
			}
		case "ss":
			validInt64, _ := strconv.ParseInt(validDate, 10, 64)
			if validInt64 < 0 || validInt64 > 59 {
				return fmt.Errorf("秒格式和值不匹配，格式：s, 值：%s", validDate)
			}
		}
		newFormat = newFormat[0:index] + newFormat[index+lenSingle:]
		newDate = newDate[0:index] + newDate[index+lenSingle:]
	}
	*date = newDate
	*format = newFormat
	return nil
}

// 验证参数
// args 参数
// minNum,maxNum 最小和最大参数数量
// ins 需要在这个列表中
func validArgs(args []interface{}, minNum, maxNum int, ins ...interface{}) error {
	if minNum < 0 {
		return fmt.Errorf("minNum必须大于等于0")
	}
	if minNum == maxNum && len(args) != minNum {
		return fmt.Errorf("验证规则错误: 参数数量必须是%d", minNum)
	}
	if len(args) < minNum || (maxNum != -1 && len(args) > maxNum) {
		if maxNum == -1 {
			return fmt.Errorf("验证规则错误: 参数数量必须大于%d", minNum)
		}
		return fmt.Errorf("验证规则错误: 参数数量必须在%d-%d之间", minNum, maxNum)
	}
	for index, arg := range args {
		if len(ins) <= index {
			continue
		}
		if ins[index] == nil {
			continue
		}
		if !inArray(arg, ins[index]) {
			return fmt.Errorf("验证规则错误: 第%d个参数%v不在%v中", index+1, arg, ins[index])
		}
	}
	return nil
}

// 公式比较
func formulaCompare(d *Data, args ...interface{}) (bool, error) {
	var formulaList [][]interface{}
	var formulaArgs []interface{}
	var formulaSymbolList []string
	for _, arg := range args {
		switch arg.(type) {
		case *condition.Formula:
			bl, err := formulaCompare(d, arg.(*condition.Formula).Args...)
			if err != nil {
				return false, err
			}
			formulaArgs = append(formulaArgs, bl)
		case string:
			switch arg.(string) {
			case "&&", "||":
				formulaSymbolList = append(formulaSymbolList, arg.(string))
				formulaList = append(formulaList, formulaArgs)
				formulaArgs = []interface{}{}
			default:
				formulaArgs = append(formulaArgs, arg)
			}
		default:
			formulaArgs = append(formulaArgs, arg)
		}
	}
	if len(formulaArgs) > 0 {
		formulaList = append(formulaList, formulaArgs)
	}
	var index int
	var isBool bool
	for index, formulaArgs = range formulaList {
		if len(formulaArgs) == 1 {
			if rs, ok := formulaArgs[0].(bool); ok {
				isBool = rs
				continue
			}
			return false, fmt.Errorf("公式错误: %v", formulaArgs)
		}
		var leftData interface{}
		var rightData interface{}
		var symbol string
		if len(formulaArgs) == 2 {
			symbol, _ = formulaArgs[0].(string)
			leftData = d.GetValidData()
			rightData = formulaArgs[1]
		} else if len(formulaArgs) == 3 {
			symbol, _ = formulaArgs[1].(string)
			leftData = formulaArgs[0]
			rightData = formulaArgs[2]
			if field, ok := formulaArgs[0].(param.Field); ok {
				rightData = true
				for _, data := range d.GetLevelData(string(field)) {
					bl, err := formulaCompare(d, data.Data, symbol, formulaArgs[2])
					if err != nil {
						return false, err
					}
					if !bl {
						rightData = false
					}
				}
				leftData = true
				symbol = "="
			}
		} else {
			return false, fmt.Errorf("公式错误: %v", formulaArgs)
		}
		formulaSymbol := ""
		if len(formulaSymbolList) > index-1 && index > 0 {
			formulaSymbol = formulaSymbolList[index-1]
		}
		var err error
		var compareData int
		var compareBool bool
		switch symbol {
		case ">":
			compareData, err = isCompareData(leftData, rightData)
			if err != nil {
				return false, err
			}
			if compareData == 1 {
				compareBool = true
			}
		case ">=":
			compareData, err = isCompareData(leftData, rightData)
			if err != nil {
				return false, err
			}
			if compareData != -1 {
				compareBool = true
			}
		case "<":
			compareData, err = isCompareData(leftData, rightData)
			if err != nil {
				return false, err
			}
			if compareData == -1 {
				compareBool = true
			}
		case "<=":
			compareData, err = isCompareData(leftData, rightData)
			if err != nil {
				return false, err
			}
			if compareData != 1 {
				compareBool = true
			}
		case "=":
			if isEqualData(leftData, rightData) {
				compareBool = true
			}
		case "!=":
			if !isEqualData(leftData, rightData) {
				compareBool = true
			}
		case "in":
			if inArray(leftData, rightData) {
				compareBool = true
			}
		case "not":
			if !inArray(leftData, rightData) {
				compareBool = true
			}
		default:
			return false, fmt.Errorf("公式错误: %v", formulaArgs)
		}
		switch formulaSymbol {
		case "&&":
			if !isBool {
				continue
			}
			isBool = compareBool
		case "||":
			if isBool {
				continue
			}
			isBool = compareBool
		default:
			isBool = compareBool
		}
	}
	return isBool, nil
}

func inArray(val interface{}, array interface{}) (exists bool) {
	exists = false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
				return
			}
		}
	}
	return
}

func inArrayRuleTree(v string, list []ruleTree) bool {
	for _, val := range list {
		if val.path == v {
			return true
		}
	}
	return false
}

// 获取float64类型数据
func toFloat64(i interface{}) float64 {
	if i == nil {
		return 0
	}
	var float64I float64
	switch i.(type) {
	case int:
		float64I = float64(i.(int))
	case int8:
		float64I = float64(i.(int8))
	case int16:
		float64I = float64(i.(int16))
	case int32:
		float64I = float64(i.(int32))
	case int64:
		float64I = float64(i.(int64))
	case uint:
		float64I = float64(i.(uint))
	case uint8:
		float64I = float64(i.(uint8))
	case uint16:
		float64I = float64(i.(uint16))
	case uint32:
		float64I = float64(i.(uint32))
	case uint64:
		float64I = float64(i.(uint64))
	case float32:
		float64I = float64(i.(float32))
	case float64:
		float64I = i.(float64)
	case string:
		float64I, _ = strconv.ParseFloat(i.(string), 64)
	case bool:
		if i.(bool) {
			float64I = 1
		}
	}
	return float64I
}

func toString(i interface{}) string {
	if i == nil {
		return ""
	}
	switch i.(type) {
	case string:
		return i.(string)
	case *string:
		return *i.(*string)
	case int, int8, int16, int32, int64:
		var i64 int64
		switch i.(type) {
		case int:
			i64 = int64(i.(int))
		case int8:
			i64 = int64(i.(int8))
		case int16:
			i64 = int64(i.(int16))
		case int32:
			i64 = int64(i.(int32))
		case int64:
			i64 = i.(int64)
		}
		return strconv.FormatInt(i64, 10)
	case *int, *int8, *int16, *int32, *int64:
		var i64 int64
		switch i.(type) {
		case *int:
			i64 = int64(*i.(*int))
		case *int8:
			i64 = int64(*i.(*int8))
		case *int16:
			i64 = int64(*i.(*int16))
		case *int32:
			i64 = int64(*i.(*int32))
		case *int64:
			i64 = *i.(*int64)
		}
		return strconv.FormatInt(i64, 10)
	case *uint, *uint8, *uint16, *uint32, *uint64:
		var ui64 uint64
		switch i.(type) {
		case *uint:
			ui64 = uint64(*i.(*uint))
		case *uint8:
			ui64 = uint64(*i.(*uint8))
		case *uint16:
			ui64 = uint64(*i.(*uint16))
		case *uint32:
			ui64 = uint64(*i.(*uint32))
		case *uint64:
			ui64 = *i.(*uint64)
		}
		return strconv.FormatUint(ui64, 10)
	case uint, uint8, uint16, uint32, uint64:
		var ui64 uint64
		switch i.(type) {
		case uint:
			ui64 = uint64(i.(uint))
		case uint8:
			ui64 = uint64(i.(uint8))
		case uint16:
			ui64 = uint64(i.(uint16))
		case uint32:
			ui64 = uint64(i.(uint32))
		case uint64:
			ui64 = i.(uint64)
		}
		return strconv.FormatUint(ui64, 10)
	case float32, float64:
		var f64 float64
		switch i.(type) {
		case float32:
			f64 = float64(i.(float32))
		case float64:
			f64 = i.(float64)
		}
		return strconv.FormatFloat(f64, 'G', -1, 64)
	case *float32, *float64:
		var f64 float64
		switch i.(type) {
		case *float32:
			f64 = float64(*i.(*float32))
		case *float64:
			f64 = *i.(*float64)
		}
		return strconv.FormatFloat(f64, 'G', -1, 64)
	case bool:
		return strconv.FormatBool(i.(bool))
	case complex64, complex128:
		var c128 complex128
		switch i.(type) {
		case complex64:
			c128 = complex128(i.(complex64))
		case complex128:
			c128 = i.(complex128)
		}
		return strconv.FormatComplex(c128, 'G', -1, 128)
	case error:
		return i.(error).Error()
	}
	return fmt.Sprintf("%v", i)
}
