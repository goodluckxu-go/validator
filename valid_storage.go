package validator

import (
	"net/http"
)

// SetRequest 设置request
func (v *valid) SetRequest(req *http.Request) (va *valid) {
	va = v.getInstance()
	va.storage.req = req
	return
}

// SetData 设置数据
func (v *valid) SetData(data interface{}) *valid {
	v.storage.data = data
	return v
}

// SetRules 设置规则
func (v *valid) SetRules(rules []Rule) (va *valid) {
	va = v.getInstance()
	va.storage.rules = rules
	return
}

// SetMessages 设置消息
func (v *valid) SetMessages(messages []Message) (va *valid) {
	va = v.getInstance()
	va.storage.messages = messages
	return
}
