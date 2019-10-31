# go-dict

> 根据结构体自动填充 string 字段

数据库中一些用字符串保存的状态，可以使用 go-dict 建立字典以访问字典属性的方式调用字符串。

```go
import "github.com/og/x/dict"

// 根据业务定义字典结构
type dictOrderStruct struct {
	Status struct{
		SystemReject string		// 默认转换为小驼峰,建议明确定义 dict ，不要使用默认转换，默认转换只是防御措施
		CheckPending string `dict:"check_pending"` // 通过 tag 自定义 value
	}
}
// 定义一个空结构体
var dictOrder = dictOrderStruct{}
// init 时通过gdict.Fill 填充结构体 (注意一定要传入指针 &dictOrder )
func init () { gdict.Fill(&dictOrder) }
func DictOrder() dictOrderStruct {
	return dictOrder
}
DictOrder().Status.SystemReject // systemReject
DictOrder().Status.CheckPending // check_pending
```

如果不想使用 godict 也可以自己手写填充，但是这样容易写错导致 key 与 value 不一致。所以使用 gdict.Fill 可自动化的赋值。

```go
type dictOrderStruct struct {
	Status struct{
		SystemReject string		// 默认转换为小驼峰
		CheckPending string `dict:"check_pending"` // 通过 tag 自定义 value
	}
}
var dictOrder = dictOrderStruct {
   Status: struct {
       SystemReject string
       CheckPending string
   }{
       SystemReject: "systemReject",
       CheckPending: "check_pending",
   },
}
func DictOrder() dictOrderStruct {
	return dictOrder
}
DictOrder().Status.SystemReject // systemReject
DictOrder().Status.CheckPending // check_pending
```