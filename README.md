# go版本的validator通用验证器
目前支持解析ValidJson,ValidXml

## 用法(usage)
只传数据data时可以只用于赋值

只传规则rule时可以只用于验证

<span style="color:red;">注意：</span>

<span style="color:red;">1. 数据data和规则rule必须传一个以上</span>

<span style="color:red;">2. 规则a.b代表a对象里面key为b的数据，a.*.b代表a数组中对象里面key为b的数据(全局通用)</span>

引入
~~~go
import "github.com/goodluckxu-go/validator"
~~~

添加全局验证器
~~~go
validator.RegisterMethod("test", func(d *validator.Data, args ...interface{}) error {
    return nil
})
~~~
调用方法
~~~go
// data只能有map[string]interface{},[]interface{},interface{}三种类型
var data interface{}
var dataIn int
err := validator.New(c.Request).ValidJson(&data, []validator.Rule{
	{Field:""}, //空字符串代表验证最外层(全部)
    {Field: "list.*.b.*.a", Methods: validator.Method.List(
	    validator.Method.SetMethod("required"), // 常规验证
	    validator.Method.SetMethod("test"), // 添加的全局验证
		// args表示外部传入的任意参数
        validator.Method.SetFun(func(d *validator.Data, args ...interface{}) error {
            a, _ := args[0].(*int)
            *a = 10
            return nil
        }, &dataIn), // 自定义验证
    ), Notes: "测试"},
}, validator.Messages{ // 只会生效最后一条
    {"list.*.b.*.a.required", "必填"}, // 其中*代表list1列表的所有的b数组中所有的a的required规则错误注释被替换
    {"list.*.b.0.a.required", "必填1"},// 其中*代表list1列表的所有的b数组中第0位的a的required规则错误注释被替换
})
~~~

## 验证数据*Data说明
### 数据
~~~json
{
    "c": "111",
    "list": [
        {
            "a": {
                "a": 1,
                "b": 2,
                "c": 3
            },
            "b": [
                {
                    "a": 5,
                    "b": 6
                },
                {
                    "a": 5,
                    "b": 7
                },
                {
                    "a": 5,
                    "b": 8
                }
            ]
        },
        {
            "a": {
                "a": 1,
                "b": 2,
                "c": 3
            },
            "b": [
                {
                    "a": 5,
                    "b": 16
                },
                {
                    "a": 5,
                    "b": 17
                },
                {
                    "a": 5,
                    "b": 18
                }
            ]
        }
    ]
}
~~~
### 方法
#### GetAllData() interface{} 
获取所有数据
#### GetData(key string) []interface{} 
根据key获取数据

例如：传list.*.b.*.b，获得[6,7,8,16,17,18]
#### GetCommonData(key string) interface{}
获取两个值最近的公共数据

例如：当期验证数据list.0.a.a，传入list.*.a.b，获取数据则为list.0.a对象
#### GetLevelData(key string) []interface{}
获取和验证同一层级的数据集合

例如：当期验证数据list.0.a.a，传入list.*.a.b，获取数据则为list.0.a.b数据切片
#### GetValidData() interface{}
获取验证数据
#### GetNotes() string
获取注释

## 返回(callback)
*valid
### 方法
Error() string   单条错误信息
Errors() []string  多条错误信息

## 语言包(language)
~~~go
validator.SetLangAddr("./zh_cn.json")
~~~
语言包json格式参照lang/zh_cn.json文件

## 所有常规验证规则(rules)
### 通用规则
errors|required

## 规则注释(rule notes)
### errors
同时验证单条数据的所有规则返回所有错误信息，设置在第一个规则，不填则验证第一个规则
### required
验证数据必填