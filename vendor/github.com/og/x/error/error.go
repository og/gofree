package ge

import (
	"fmt"
	gdict "github.com/og/x/dict"
	"runtime/debug"
	"strings"
)
func FillErrCode(v interface{}) {
	gdict.CustomFill(v, gdict.Custom{
		StructTagName: "dict",
		ValueASKey: true,
	})
}

type Err struct {
	fail bool
	msg string
	code string
	stack []byte
}
func (err Err) Fail() bool {
	return err.fail
}
func (err Err) Code() string {
	return err.code
}
func (err Err) Msg() string {
	return err.msg
}
func byteIsCapital(b byte) bool {
	s := string(b)
	return s <= "Z" && s >= "A"
}
func (err *Err) SetCode(code string) {
	err.stack = debug.Stack()
	err.fail = true
	err.code = code
	var msgByteList []byte
	for i:=0;i<len(code);i++ {
		item := code[i]
		if byteIsCapital(item){
			if i !=0 && i != len(code)-1 && (!byteIsCapital(code[i-1]) || !byteIsCapital(code[i+1])) {
				msgByteList = append(msgByteList, " "[0])
			}
			if i!=0 {
				item = strings.ToLower(string(item))[0]
			}
		}
		msgByteList = append(msgByteList, item)
	}
	err.msg = string(msgByteList[:])
}
func (err *Err) SetCodeMsg(code string, msg string) {
	err.stack = debug.Stack()
	err.fail = true
	err.code = code
	err.msg = msg
}
func (err Err) Error() string {
	return fmt.Sprintf("code:\r\n\t%s\r\nnmessage:\r\n\t%s\r\n%s", err.code, err.msg, err.stack)
}