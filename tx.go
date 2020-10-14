package f

import (
	"github.com/jmoiron/sqlx"
	ge "github.com/og/x/error"
	"log"
)

func newTx(tx *sqlx.Tx) *Tx {
	return &Tx{core: tx}
}
type Tx struct {
	done bool
	core *sqlx.Tx
}
func (tx *Tx) Commit() {
	if tx.done {
		log.Print("Many times to Commit")
	} else {
		ge.Check(tx.core.Commit())
		tx.done = true
	}
}

func (tx *Tx) Rollback() {
	if tx.done {
		log.Print("Many times to Rollback")
	} else {
		ge.Check(tx.core.Rollback())
		tx.done = true
	}
}
