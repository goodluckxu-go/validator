# go版本的validator通用验证器

## 修改日志(change log)
改版前v1.0.1

![](./改版前.png)

改版后v1.0.4

![](./改版后.png)

树结构v1.0.5

![](./树结构.png)

## 用法(usage)

<span style="color:red">注意：</span>

<span style="color:red">1. 规则a.b代表a对象里面key为b的数据，a.*.b代表a数组中对象里面key为b的数据(全局通用)</span>

<span style="color:red">2. 验证参数存在 param.Field("fieldName") 时，fieldName 必须存在验证规则中</span>

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
// data只能有map[string]interface{},[]interface{},interface{}三种类型时字符串类型的数字可转化
var data interface{}
var dataIn int
valid := validator.New().
	SetRequest(req). // *http.Request，不设置则验证传入data数据的值
	SetData(&data).
	SetRules([]validator.Rule{
		{Field:""}, //空字符串代表验证最外层(全部) 
		{Field: "list.*.b.*.a", Methods: validator.Method.List(
			validator.Method.SetMethod("required"), // 常规验证 
			validator.Method.SetMethod("test"), // 添加的全局验证 
			//args表示外部传入的任意参数 
			validator.Method.SetFun(func(d *validator.Data, args ...interface{}) error {
				a, _ := args[0].(*int)
				*a = 10
				return nil
			}, &dataIn), // 自定义验证 
		), Notes: "测试"},
	}). // 不传规则则是赋值
	SetMessages([]validator.Message{ // 只会生效最后一条 
		{"list.*.b.*.a.required", "必填"}, // 其中*代表list1列表的所有的b数组中所有的a的required规则错误注释被替换 
		{"list.*.b.0.a.required", "必填1"},// 其中*代表list1列表的所有的b数组中第0位的a的required规则错误注释被替换 
	}).
	Valid()
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
#### GetData(path string) []DataOne
根据key获取数据

例如：传 list.*.b. *.b ，获得
~~~json
[
    {"Path":"list.0.b.0.b","Data":6},
    {"Path":"list.0.b.1.b","Data":7},
    {"Path":"list.0.b.2.b","Data":8},
    {"Path":"list.1.b.0.b","Data":16},
    {"Path":"list.1.b.1.b","Data":17},
    {"Path":"list.1.b.2.b","Data":18}
]
~~~
#### GetCommonData(path string) interface{}
获取两个值最近的公共数据

例如：当期验证数据list.0.a.a，传入list.*.a.b，获取数据则为list.0.a对象
#### GetLevelData(path string) []DataOne
获取和验证同一层级的数据集合

例如：当期验证数据list.0.a.a，传入list.*.a.b，获取数据则为list.0.a.b数据切片
#### GetValidData() interface{}
设置验证数据
#### GetNotes() string
获取注释
#### GetNotesByPath(path string) string
通过path路径获取注释
#### func (d *Data) JumpValid() error
跳过当前字段的验证
#### func (d *Data) JumpChild() error {
跳过子集验证

## 返回错误信息(callback)
valid.Error

## 语言包(language)
~~~go
validator.SetLangAddr("./zh_cn.json")
~~~
语言包json格式参照lang/zh_cn.json文件

## 所有常规验证规则(rules)
### 通用规则
[valid_condition](#valid_condition) |
[required](#required) |
[nullable](#nullable) |
[in](#in) |
[not_in](#not_in) |
[unique](#unique) |
[email](#email) |
[phone](#phone) |
[regexp](#regexp) |
[not_regexp](#not_regexp)
### 类型验证
[array](#array) |
[map](#map) |
[string](#string) |
[number](#number) |
[integer](#integer) |
[bool](#bool) |
[date](#date) |
[file](#file)
### 比较验证 (字段,时间,数字)
[eq](#eq) |
[gt](#gt) |
[gte](#gte) |
[lt](#lt) |
[lte](#lte)
### 长度验证 (字符串,数组,文件大小)
[len](#len) |
[min](#min) |
[max](#max)
### 文件验证
[suffix](#suffix) |
[mime](#mime)

## 规则注释(rule notes)

#### <a id="valid_condition">valid_condition规则</a>
<span style="color:red;">args参数说明:</span> 

<span style="color:red;">1. &&是并且,||是或者</span>

<span style="color:red;">2. &&和||之间的数据为验证公式，公式结果为布尔型</span>

<span style="color:red;">3. 公式目前占1-3个字符</span>

<span style="color:red;">4. 公式可以直接为true或false</span>

<span style="color:red;">5. 验证公式符号有>,>=,<,<=,=,!=,in,not。例如：>,3表示验证数据大于3；3,>,2表示比较3和2的值；param.Field("a"),>,2表示字段a的值大于2</span>

<span style="color:red;">6. in是在数组里面，not是不再数组里面，其他公式符号都是简单类型数组(string,int等)</span>

<span style="color:red;">7. condition.Brackets表示括号，里面值和args一样规则</span>

<span style="color:red;">8. 所有公式满足为true则验证，否在跳过验证</span>
~~~go
validator.Method.SetMethod("valid_condition", "=", 1, "&&", "<", 10, "&&", param.Field("list.*.a.a"), "=", 12, "&&", condition.Brackets(">", 2)),
~~~

#### <a id="required">required规则</a>
验证是否必填。null，字符串为""，数字类型为0，bool类型为false，数组为[]，map为{}都不通过
参数为string,number,bool,array,map，代表string为空不验证，number为0不验证，bool为false不验证,array为[]不验证,map为{}不验证
~~~go
validator.Method.SetMethod("required", "string")
~~~

#### <a id="nullable">nullable规则</a>
数据为空则不验证后面的规则，字符串为""，数字类型为0，bool类型为false，数组为[]，map为{}都不验证后的单规则
~~~go
validator.Method.SetMethod("nullable")
~~~

#### <a id="in">in规则</a>
验证是否在数组里面，参数是slice或array类型
~~~go
validator.Method.SetMethod("in", []int{1,2,3})
~~~

#### <a id="not_in">not_in规则</a>
验证是否不在数组里面，参数是slice或array类型
~~~go
validator.Method.SetMethod("not_in", []int{1,2,3})
~~~

#### <a id="unique">unique规则</a>
验证数组内的值唯一
~~~go
validator.Method.SetMethod("unique")
~~~

#### <a id="email">email规则</a>
验证是否是邮箱
~~~go
validator.Method.SetMethod("email")
~~~

#### <a id="phone">phone规则</a>
验证是否是手机号
~~~go
validator.Method.SetMethod("phone")
~~~

#### <a id="regexp">regexp规则</a>
验证数据在正则表达式中
~~~go
validator.Method.SetMethod("regexp", `^\d*$`)
~~~

#### <a id="not_regexp">not_regexp规则</a>
验证数据不在正则表达式中
~~~go
validator.Method.SetMethod("not_regexp", `^\d*$`)
~~~

#### <a id="array">array规则</a>
验证是否是数组
~~~go
validator.Method.SetMethod("array")
~~~

#### <a id="map">map规则</a>
验证是否是对象
~~~go
validator.Method.SetMethod("map")
~~~

#### <a id="string">string规则</a>
验证是否是字符串
~~~go
validator.Method.SetMethod("string")
~~~

#### <a id="number">number规则</a>
验证是否是数字。可验证数字和字符串的数字，如果接受数据为map[string]interface{},[]interface{},interface{}验证后会转换成float64
~~~go
validator.Method.SetMethod("number")
~~~

#### <a id="integer">integer规则</a>
验证是否是整数。可验证数字和字符串的数字，如果接受数据为map[string]interface{},[]interface{},interface{}验证后会转换成int64
~~~go
validator.Method.SetMethod("integer")
~~~

#### <a id="bool">bool规则</a>
验证是否是布尔类型。可验证整数和布尔类型
~~~go
validator.Method.SetMethod("bool")
~~~

#### <a id="date">date规则</a>
验证是否是日期格式 args参数Y-m-d H:i:s类型，Y年，m月，d日，H时，i分，s秒
~~~go
validator.Method.SetMethod("date")
~~~

#### <a id="file">file规则</a>
验证是否是文件
~~~go
validator.Method.SetMethod("file")
~~~

#### <a id="eq">eq规则</a>
验证两个字段是否相同，参数可以是param.Field("test")字段值
~~~go
validator.Method.SetMethod("eq",5)
~~~

#### <a id="gt">gt规则</a>
验证是否大于某个数。可验证数字，字符串的数字或日期，参数可以是param.Field("test"),time.Time字段值
~~~go
validator.Method.SetMethod("gt",5)
~~~

#### <a id="gte">gte规则</a>
验证是否大于等于某个数。可验证数字和字符串的数字或日期，参数可以是param.Field("test"),time.Time字段值
~~~go
validator.Method.SetMethod("gte",5)
~~~

#### <a id="lt">lt规则</a>
验证是否小于某个数。可验证数字和字符串的数字或日期，参数可以是param.Field("test"),time.Time字段值
~~~go
validator.Method.SetMethod("lt",5)
~~~

#### <a id="lte">lte规则</a>
验证是否小于等于某个数。可验证数字和字符串的数字或日期，参数可以是param.Field("test"),time.Time字段值
~~~go
validator.Method.SetMethod("lte",5)
~~~

#### <a id="len">len规则</a>
验证是长度等于某个数。可获取字符串和数组长度或者文件大小(单位: kb)
~~~go
validator.Method.SetMethod("len", 5)
~~~

#### <a id="min">min规则</a>
验证是长度大于等于某个数。可获取字符串和数组长度或者文件大小(单位: kb)
~~~go
validator.Method.SetMethod("min", 5)
~~~

#### <a id="max">max规则</a>
验证是长度小于等于某个数。可获取字符串和数组长度或者文件大小(单位: kb)
~~~go
validator.Method.SetMethod("max", 5)
~~~

#### <a id="suffix">suffix规则</a>
验证文件后缀类型，可为多个
~~~go
validator.Method.SetMethod("suffix", "jpg","png")
~~~

#### <a id="mime">mime规则</a>
验证文件协议类型，可为多个
~~~go
validator.Method.SetMethod("mime", "image/jpeg","image/png")
~~~