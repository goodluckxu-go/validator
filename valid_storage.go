package validator

import (
	"net/http"
)

// SetRequest 设置request
func (v *Valid) SetRequest(req *http.Request) (va *Valid) {
	va = v.getInstance()
	va.storage.req = req
	return
}

// SetData 设置数据
func (v *Valid) SetData(data interface{}) (va *Valid) {
	va = v.getInstance()
	va.storage.data = data
	return
}

// SetRules 设置规则
func (v *Valid) SetRules(rules []Rule) (va *Valid) {
	va = v.getInstance()
	va.storage.rules = rules
	return
}

// SetMessages 设置消息
func (v *Valid) SetMessages(messages []Message) (va *Valid) {
	va = v.getInstance()
	va.storage.messages = messages
	return
}
