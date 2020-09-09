package f_test

import (
	f "github.com/og/gofree"
	ge "github.com/og/x/error"
	"testing"
)

var mock = f.Mock{
	Tables:[]interface{}{
		[]User{
			{Name:"nimo"},
			{Name:"nico"},
		},
	},
}
func TestMock(t *testing.T) {
	db, err := NewDB() ; ge.Check(err)
	f.ResetAndMock(db, mock)
}
