package f

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

func newTx(tx *sqlx.Tx) *Tx {
	return &Tx{Core: tx}
}
type Tx struct {
	Core *sqlx.Tx
}

func (tx *Tx) rollback() error {
	err := tx.Core.Rollback()
	if errors.Is(err, sql.ErrTxDone) {
		// 忽略 tx done
	} else {
		return err
	}
	return nil
}
func (tx *Tx) commit() error {
	err := tx.Core.Commit()
	if errors.Is(err, sql.ErrTxDone) {
		// 忽略 tx done
	} else {
		return err
	}
	return nil
}