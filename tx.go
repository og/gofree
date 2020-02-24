package f

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	ge "github.com/og/x/error"
)
// 将 tx.Rollback() 返回的错误 panic(err)
func Rollback(tx *sqlx.Tx) {
	ge.Check(tx.Rollback())
}
// 自动提交事务
// 在使用 `tx := db.Core.MustBegin()` 下一行加上 `defer func() { f.AutoCommitTx(tx, recover()) }()`
// 能够在函数执行结束时判断如何没有任何 `recover` 并没有执行过 `tx.Commit()` 或 `tx.Rollback()` 时执行 `tx.Commit()`
// 这样只需要在逻辑代码中需取消事务的地方加上 tx.Rollback() 或者 f.Rollback()
func AutoCommit( tx *sqlx.Tx, recover interface {}) {
	if recover == nil {
		err := tx.Commit()
		if err == sql.ErrTxDone{
			// break
		} else {
			ge.Check(err)
		}
	} else {
		err := tx.Rollback()
		if err == sql.ErrTxDone{
			panic(recover)
		} else {
			panic(err)
		}
	}
}