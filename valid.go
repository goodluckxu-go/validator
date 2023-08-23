package validator

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
)

type Valid struct {
	storage storage // 存储仓库
	handle  handle  // 处理数据
	Error   error
}

// Valid 验证数据
func (v *Valid) Valid() (va *Valid) {
	defer func() {
		r := recover()
		if r != nil {
			va.Error = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	va = v.getInstance()
	var err error
	if va.storage.req != nil {
		err = va.parseRequest()
		if err != nil {
			va.Error = err
			return
		}
	}
	if va.storage.rules != nil {
		dataValue := reflect.ValueOf(va.storage.data).Elem()
		newData := dataValue.Interface()
		isSet := false
		switch newData.(type) {
		case map[string]interface{}, []interface{}:
			isSet = true
		default:
			by, byErr := json.Marshal(newData)
			if byErr != nil {
				va.Error = byErr
				return
			}
			if err = json.Unmarshal(by, &newData); err != nil {
				va.Error = err
				return
			}
		}
		if err = v.parseRule(va.storage.rules, newData); err != nil {
			va.Error = err
			return
		}
		if err = v.validRule(&newData); err != nil {
			va.Error = err
			return
		}
		if isSet {
			dataValue.Set(reflect.ValueOf(newData))
		}
	}
	return
}

// 重新实例化
func (v *Valid) getInstance() (va *Valid) {
	va = new(Valid)
	va.storage = v.storage
	va.handle = v.handle
	return
}

// parseRequest 解析请求数据
func (v *Valid) parseRequest() error {
	if v.storage.req == nil {
		return nil
	}
	contentTypeList := strings.Split(v.storage.req.Header.Get("Content-Type"), ";")
	contentType := contentTypeList[0]
	switch contentType {
	case "application/json":
		body := readBody(v.storage.req)
		if err := json.Unmarshal(body, v.storage.data); err != nil {
			return err
		}
	case "application/xml":
		body := readBody(v.storage.req)
		if err := xml.Unmarshal(body, v.storage.data); err != nil {
			return err
		}
	case "application/x-www-form-urlencoded":
		body := readBody(v.storage.req)
		bodyList := strings.Split(string(body), "&")
		data := map[string]interface{}{}
		for _, oneBody := range bodyList {
			oneList := strings.Split(oneBody, "=")
			oneKey, _ := url.PathUnescape(oneList[0])
			oneVal, _ := url.PathUnescape(oneList[1])
			data[oneKey] = oneVal
		}
		reflect.ValueOf(v.storage.data).Elem().Set(reflect.ValueOf(data))
	case "multipart/form-data":
		if len(contentTypeList) < 2 {
			return errors.New("请传入 multipart/form-data 格式数据")
		}
		boundary := strings.Split(contentTypeList[1], "=")[1]
		body := readBody(v.storage.req)
		bodyList := strings.Split(string(body), boundary)
		bodyList = bodyList[1 : len(bodyList)-1]
		data := map[string]interface{}{}
		for _, oneBody := range bodyList {
			oneList := strings.Split(oneBody, "\r\n")
			explanList := strings.Split(oneList[1], ";")
			name := strings.Trim(strings.Split(explanList[1], "=")[1], "\"")
			if len(explanList) == 2 {
				// 常规数据
				data[name] = oneList[3]
			} else {
				// 文件数据
				if _, fh, err := v.storage.req.FormFile(name); err == nil {
					v.handle.fileMap[name] = &file{
						Suffix: strings.TrimPrefix(filepath.Ext(fh.Filename), "."),
						Mime:   fh.Header.Get("Content-Type"),
						Name:   fh.Filename,
						Size:   fh.Size,
					}
				}
			}
		}
		reflect.ValueOf(v.storage.data).Elem().Set(reflect.ValueOf(data))
	default:
		body := readBody(v.storage.req)
		reflect.ValueOf(v.storage.data).Elem().Set(reflect.ValueOf(string(body)))
	}
	return nil
}

// parseRule 解析规则
func (v *Valid) parseRule(rules []Rule, data interface{}) error {
	ruleRowList, err := disintegrateRules(rules, data, true)
	if err != nil {
		return err
	}
	v.handle.ruleRowList, v.handle.pathIndex = ruleRowSort(ruleRowList)
	return nil
}

// validRule 验证规则
func (v *Valid) validRule(data *interface{}) (err error) {
	for _, ruleOnce := range v.handle.ruleRowList {
		// 验证数据
		validNum := 0
		for _, m := range ruleOnce.methods {
			validNum++
			if ruleOnce.notes == "" {
				ruleOnce.notes = ruleOnce.path
			}
			d := &Data{
				data:           data,
				path:           ruleOnce.path,
				ruleRowListPtr: &v.handle.ruleRowList,
				pathIndexPtr:   &v.handle.pathIndex,
				fileMapPtr:     &v.handle.fileMap,
				messagesPtr:    &v.storage.messages,
			}
			var fn methodFunc
			var me string
			switch m.method.(type) {
			case string:
				me = m.method.(string)
				d.methodName = me
				var fnInterface interface{}
				var ok bool
				if fnInterface, ok = methodPool.Load(me); !ok {
					err = fmt.Errorf("规则 %s 不存在", me)
					return
				}
				if fn, ok = fnInterface.(methodFunc); !ok {
					err = errors.New("规则没有注入")
					return
				}

			case methodFunc:
				fn = m.method.(methodFunc)
			default:
				err = errors.New("未知错误")
				return
			}
			err = fn(d, m.args...)
			if err == nil {
				continue
			}
			if err.Error() == jumpValid {
				err = nil
				break
			}
			return
		}
		if len(ruleOnce.methods) == validNum {
			v.handle.ruleRowList[v.handle.pathIndex[ruleOnce.path]].isValid = true
		}
	}
	return
}
