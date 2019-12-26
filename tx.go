package f

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	ge "github.com/og/x/error"
)

func EndTx( tx *sqlx.Tx, recover interface {}) {
	if recover == nil {
		err := tx.Commit()
		if err == sql.ErrTxDone && err == nil {
			// break
		} else {
			ge.Check(err)
		}
	} else {
		err := tx.Rollback()
		if err == sql.ErrTxDone && err == nil {
			panic(recover)
		} else {
			panic(err)
		}
	}
}