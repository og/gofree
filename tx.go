package f

import (
	"github.com/jmoiron/sqlx"
	ge "github.com/og/x/error"
)

func newTx(tx *sqlx.Tx) Tx {
	return Tx{core: tx}
}
type Tx struct {
	done bool
	core *sqlx.Tx
}
func (tx Tx) Commit() {
	if !tx.done {
		ge.Check(tx.core.Commit())
		tx.done = true
	}
}

func (tx Tx) Rollback() {
	if !tx.done {
		ge.Check(tx.core.Rollback())
		tx.done = true
	}
}
// 使用此函数千万注意不要出现多个 defer 都运行 recover() ，并且 recover 不是 nil 时会向上传递
// 建议阅读此函数源码了解运行机制（代码很简单的哦）
func (tx Tx) End(recoverValue interface{}) {
	if recoverValue == nil {
		tx.Commit()
	} else {
		tx.Rollback()
		panic(recoverValue) // 将错误传递
	}
}