package scan_test

import (
	"github.com/jmoiron/sqlx"
	"github.com/og/gofree/scan"
	ge "github.com/og/x/error"
	gjson "github.com/og/x/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapScanStruct(t *testing.T) {
		db, err := sqlx.Open("mysql","root:password@(localhost:3306)/test_gofree?charset=utf8&loc=Local&parseTime=True") ; if err != nil {panic(err)}
		row, err := db.Queryx(`SELECT
		  user.id             AS "user.id"
		, user.name           AS "user.name"
		, user.is_super       AS "user.is_super"
	
		, book.id             AS "book.id"
		, book.name           AS "book.name"
		, book.price          AS "book.price"
		, book.user_id        AS "book.user_id"
		, book.last_read_time AS "book.last_read_time"
	
	FROM
		user JOIN book ON user.id = book.user_id
	WHERE
		user.id = ?`, 1) ; ge.Check(err)
		userAndBook := UserAndBook{}
		userAndBookScan := scan.New(&userAndBook)
		mapRelation := scan.GetRelation(userAndBook)
		for row.Next() {
			data := map[string]interface{}{}
			ge.Check(row.MapScan(data))
			userAndBookScan.MapScanStruct(data, mapRelation, userAndBook)
		}
		assert.Equal(t, `{
  "User": {
    "ID": 1,
    "Name": "bmltbw==",
    "IsSuper": 0
  },
  "BookList": [
    {
      "ID": 1,
      "Name": "Z28gYWN0aW9u",
      "Price": "MTIuMDA=",
      "UserID": 1,
      "LastReadTime": "2019-11-11T00:00:00+08:00"
    },
    {
      "ID": 2,
      "Name": "anMgYWN0aW9u",
      "Price": "MTEuMDA=",
      "UserID": 1,
      "LastReadTime": "2019-11-11T00:00:00+08:00"
    }
  ]
}`, gjson.StringUnfold(userAndBook))
}