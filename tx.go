package f

import (
	"github.com/jmoiron/sqlx"
)

func newTx(tx *sqlx.Tx) *Tx {
	return &Tx{Tx: tx}
}
type Tx struct {
	done bool
	*sqlx.Tx
}
