package f

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

func newTransaction(tx *sqlx.Tx) *Transaction {
	return &Transaction{Core: tx}
}
type Transaction struct {
	Core *sqlx.Tx
}

func (tx *Transaction) Rollback() error {
	err := tx.Core.Rollback()
	if errors.Is(err, sql.ErrTxDone) {
		// 忽略 tx done
	} else {
		return err
	}
	return nil
}
func (tx *Transaction) commit() error {
	err := tx.Core.Commit()
	if errors.Is(err, sql.ErrTxDone) {
		// 忽略 tx done
	} else {
		return err
	}
	return nil
}