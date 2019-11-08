package scan_test

import (
	"github.com/og/gofree/scan"
	ge "github.com/og/x/error"
	gjson "github.com/og/x/json"
	gtime "github.com/og/x/time"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)



func TestMapScanStruct(t *testing.T) {
	dataList := []map[string]interface{}{
		{
			"user.id": 1,
			"user.name": "nimo",
			"user.is_super": false,
			"book.id": 1,
			"book.name": "go action",
			"book.price": 11.00,
			"book.user_id": 1,
			"book.last_read_time": ge.GetTime(time.Parse(gtime.Second, "2019-11-11 00:00:00")),
		},
		{
			"user.id": 1,
			"user.name": "nimo",
			"user.is_super": false,
			"book.id": 2,
			"book.name": "js action",
			"book.price": 11.00,
			"book.user_id": 1,
			"book.last_read_time": ge.GetTime(time.Parse(gtime.Second, "2019-11-12 00:00:00")),
		},
	}
	userAndBook := UserAndBook{}
	userAndBookScan := scan.New(&userAndBook)
	mapRelation := scan.GetRelation(userAndBook)
	for i:=0;i<len(dataList);i++ {
		data := dataList[i]
		userAndBookScan.MapScanStruct(data, mapRelation)
	}
	assert.Equal(t, `{
  "User": {
    "ID": 1,
    "Name": "nimo",
    "IsSuper": false
  },
  "BookList": [
    {
      "ID": 1,
      "Name": "go action",
      "Price": 11,
      "UserID": 1,
      "LastReadTime": "2019-11-11T00:00:00Z"
    },
    {
      "ID": 2,
      "Name": "js action",
      "Price": 11,
      "UserID": 1,
      "LastReadTime": "2019-11-12T00:00:00Z"
    }
  ]
}`, gjson.StringUnfold(userAndBook))
}
func TestScan_MapScanSlice(t *testing.T) {
	userList := []map[string]interface{}{
		{
			"user.id": 1,
			"user.name": "nimo",
			"user.is_super": false,
		},
		{
			"user.id": 2,
			"user.name": "nico",
			"user.is_super": true,
		},
	}
	_=userList
	bookList := []map[string]interface{}{
		{
			"book.id": 1,
			"book.name": "go action",
			"book.price": 12,
			"book.user_id": 1,
			"book.last_read_time": ge.GetTime(time.Parse(gtime.Second, "2019-11-12 00:00:00")),
		},
		{
			"book.id": 2,
			"book.name": "js action",
			"book.price": 11,
			"book.user_id": 1,
			"book.last_read_time": ge.GetTime(time.Parse(gtime.Second, "2019-11-12 00:00:00")),
		},
		{
			"book.id": 3,
			"book.name": "life",
			"book.price": 1,
			"book.user_id": 2,
			"book.last_read_time": ge.GetTime(time.Parse(gtime.Second, "2019-11-12 00:00:00")),
		},
	}
	_=bookList
}