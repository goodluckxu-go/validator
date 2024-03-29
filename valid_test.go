package validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/goodluckxu-go/validator/param"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

var req *http.Request
var reqOnce sync.Once

func init() {
	// 全局验证
	RegisterMethod("linda_auth", func(d *Data, args ...interface{}) error {
		if d.GetValidData() == "linda" && args[0] == "admin" && args[1] == 123456 {
			return nil
		}
		return fmt.Errorf("%s不是linda", d.GetNotes())
	})
}

func getRequest() *http.Request {
	reqOnce.Do(func() {
		req, _ = http.NewRequest("GET", "/", io.NopCloser(bytes.NewBuffer(getJsonBody())))
		req.Header.Add("Content-Type", "application/json")
	})
	return req
}

func getBody() interface{} {
	return map[string]interface{}{
		"user": map[string]interface{}{
			"username": "linda",
			"age":      15,
			"birthday": "2008-01-01",
			"pwd":      "123456",
			"re_pwd":   "123456",
			"is_vip":   true,
		},
		"goods": []interface{}{
			map[string]interface{}{
				"id":        1,
				"name":      "毛巾",
				"number":    5,
				"money":     15.25,
				"type":      5,
				"explan":    "必需品",
				"is_number": true,
			},
			map[string]interface{}{
				"id":     2,
				"name":   "衣服",
				"number": 5,
				"money":  15.25,
				"type":   5,
				"explan": "",
			},
			map[string]interface{}{
				"id":     3,
				"name":   "裤子",
				"number": 5,
				"money":  15.25,
				"type":   5,
			},
		},
	}
}

func getJsonBody() []byte {
	rs, _ := json.Marshal(getBody())
	return rs
}

func getRules() []Rule {
	return []Rule{
		{Field: "user", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("map"),
		), Notes: "用户信息"},
		{Field: "user.username", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("string"),
			Method.SetMethod("linda_auth", "admin", 123456),
			Method.SetMethod("min", 1),
			Method.SetMethod("max", 10),
			Method.SetMethod("len", 5),
		), Notes: "用户名"},
		{Field: "user.age", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("integer"),
			Method.SetMethod("gt", 0),
			Method.SetMethod("lt", 100),
		), Notes: "年龄"},
		{Field: "user.birthday", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("date", "Y-m-d"),
			Method.SetMethod("gt", "2001-01-01"),
			Method.SetMethod("lt", time.Now()),
		), Notes: "生日"},
		{Field: "user.pwd", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("string"),
		), Notes: "密码"},
		{Field: "user.re_pwd", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("string"),
			Method.SetMethod("eq", param.Field("user.pwd")),
		), Notes: "确认密码"},
		{Field: "user.is_vip", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("bool"),
		)},
		{Field: "goods", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("array"),
		), Notes: "商品列表"},
		{Field: "goods.*", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("map"),
		), Notes: "商品信息"},
		{Field: "goods.*.id", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("unique"),
			Method.SetMethod("integer"),
			Method.SetMethod("gt", 0),
		)},
		{Field: "goods.*.name", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("string"),
			Method.SetMethod("max", 255),
		)},
		{Field: "goods.*.number", Methods: Method.List(
			Method.SetMethod("valid_condition", param.Field("goods.*.is_number"), "=", true),
			Method.SetMethod("required"),
			Method.SetMethod("integer"),
			Method.SetMethod("gte", 1),
			Method.SetMethod("lte", 99),
		)},
		{Field: "goods.*.money", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("number"),
			Method.SetMethod("gt", 0),
		)},
		{Field: "goods.*.type", Methods: Method.List(
			Method.SetMethod("required"),
			Method.SetMethod("number"),
			Method.SetMethod("in", []int{1, 3, 5, 7, 9}),
			Method.SetMethod("not_in", []int{15, 26}),
		)},
		{Field: "goods.*.explan", Methods: Method.List(
			Method.SetMethod("nullable"),
			Method.SetMethod("string"),
		)},
	}
}

func BenchmarkValid_Valid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := map[string]interface{}{}
		v := New().
			SetRequest(getRequest()).
			SetData(&data).
			SetRules(getRules()).
			SetMessages([]Message{
				{"goods.*.type.required", "商品的类型必须传入"},
			}).
			Valid()
		if v.Error != nil {
			b.Errorf(v.Error.Error())
		}
	}
}

func TestValid_Valid(t *testing.T) {
	data := map[string]interface{}{}
	v := New().
		SetRequest(getRequest()).
		SetData(&data).
		SetRules(getRules()).
		SetMessages([]Message{
			{"goods.*.type.required", "商品的类型必须传入"},
			{"goods.*.explan.nullable", "商品的类型必须传入1"},
			{"goods.2.explan.nullable", "商品的类型必须传入2"},
		}).
		Valid()
	if v.Error != nil {
		t.Errorf(v.Error.Error())
	}
}
