package f_test

import (
	f "github.com/og/gofree"
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
	db := NewDB()
	f.ResetAndMock(db, mock)
}
