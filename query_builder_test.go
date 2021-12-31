package f_test

import (
	f "github.com/og/gofree"
	ge "github.com/og/x/error"
	gtest "github.com/og/x/test"
	gtime "github.com/og/x/time"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)


func TestQB_Sql(t *testing.T) {
	as := gtest.NewAS(t)
	{
		qb := f.QB{
			Table: "user",
			Select: []string{"name"},
			Where: []f.AND{
				{"name": {f.Eql("nimo")}},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT `name` FROM `user` WHERE `name` = ?", sqlS)
		assert.Equal(t, []interface {}{"nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Select: []string{"name"},
			Where: []f.AND{
				{"name": {f.Eql("nimo")}},
			},
			Check: []string{"SELECT `name` FROM `user` WHERE `name` = ?"},
		}
		_, _ = qb.GetSelect()
	}
	{
		qb := f.QB{
			Table: "user",
			Select: []string{"name", "id"},
			Where: []f.AND{
				{"name": {f.Eql("nimo")}},
			},
			Check: []string{"SELECT `name`, `id` FROM `user` WHERE `name` = ?","SELECT `name` FROM `user` WHERE `name` = ?"},
		}
		_, _ = qb.GetSelect()
	}
	{
		qb := f.QB{
			Table: "user",
			Select: []string{"name", "age"},
			Where: []f.AND{
				{"name": {f.Eql("nimo")}},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT `name`, `age` FROM `user` WHERE `name` = ?", sqlS)
		assert.Equal(t, []interface {}{"nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.And("name","nimo"),
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `name` = ?", sqlS)
		assert.Equal(t, []interface {}{"nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{"name": f.OP{f.Eql("nimo")}},
			},
			SoftDelete: "deleted_at",
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `name` = ? AND `deleted_at` IS NULL", sqlS)
		assert.Equal(t, []interface {}{"nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.And("name", "nimo"),
			SoftDelete: "deleted_at",
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `name` = ? AND `deleted_at` IS NULL", sqlS)
		assert.Equal(t, []interface {}{"nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{"name": f.OP{f.Eql("nico")}},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `name` = ?", sqlS)
		assert.Equal(t, []interface {}{"nico"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.And("name", f.Eql("nico")),
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `name` = ?", sqlS)
		assert.Equal(t, []interface {}{"nico"}, values)
	}

	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"name": f.OP{f.Eql("nimo")},
					"age":  f.OP{f.Eql(18)},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `age` = ? AND `name` = ?", sqlS)
		assert.Equal(t, []interface {}{18, "nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.And("name", "nimo", "age",18),
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `age` = ? AND `name` = ?", sqlS)
		assert.Equal(t, []interface {}{18, "nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"name": f.OP{f.Eql("nimo")},
					"age": f.OP{f.Lt(18), f.Gt(19)},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `age` < ? AND `age` > ? AND `name` = ?", sqlS)
		assert.Equal(t, []interface {}{18, 19, "nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.And(
				"name", f.Eql("nimo"),
				"age", f.Lt(18),
				"age", f.Gt(19)),
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `age` < ? AND `age` > ? AND `name` = ?", sqlS)
		assert.Equal(t, []interface {}{18, 19, "nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
					{"city": f.OP{f.Eql("shanghai")}},
					{
						"name": f.OP{f.Eql("nimo")},
						"age": f.OP{f.Lt(18), f.Gt(19)},
					},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE ( `city` = ? ) OR ( `age` < ? AND `age` > ? AND `name` = ? )", sqlS)
		assert.Equal(t, []interface{}{"shanghai", 18, 19, "nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.Or(
				f.And("city", "shanghai"),
				f.And(
					"name",f.Eql("nimo"),
					"age", f.Lt(18),
					"age", f.Gt(19),
					),
				) ,
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE ( `city` = ? ) OR ( `age` < ? AND `age` > ? AND `name` = ? )", sqlS)
		assert.Equal(t, []interface{}{"shanghai", 18, 19, "nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.Or(
				f.And("city", "shanghai"),
				f.And(
					"name",f.Eql("nimo"),
					"age", f.Lt(18),
					"age", f.Gt(19),
				),
				f.And("created_at", f.Day(ge.Time(time.Parse(gtime.LayoutTime, "2018-11-11 00:11:11")))),
			) ,
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE ( `city` = ? ) OR ( `age` < ? AND `age` > ? AND `name` = ? ) OR ( `created_at` >= ? AND `created_at` <= ? )", sqlS)
		assert.Equal(t, []interface{}{"shanghai", 18, 19, "nimo",  "2018-11-11 00:00:00", "2018-11-11 23:59:59"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"name": f.OP{f.Eql("nimo")},
					"age": f.OP{f.LtEql(18), f.GtEql(19)},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `age` <= ? AND `age` >= ? AND `name` = ?", sqlS)
		assert.Equal(t, []interface {}{18, 19, "nimo"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"id": f.OP{
						f.Like("1"),
					},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `id` LIKE ?", sqlS)
		assert.Equal(t, []interface{}{"%1%"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"id": f.OP{
						f.LikeStart("1"),
					},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `id` LIKE ?", sqlS)
		assert.Equal(t, []interface{}{"1%"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"id": f.OP{
						f.LikeEnd("1"),
					},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `id` LIKE ?", sqlS)
		assert.Equal(t, []interface{}{"%1"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"id": f.OP{
						f.LikeStart("1"),
					},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `id` LIKE ?", sqlS)
		assert.Equal(t, []interface{}{"1%"}, values)
	}
	{
		idList := []int{4,5,6}
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"id": f.OP{
						f.In(idList),
					},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `id` IN (?, ?, ?)", sqlS)
		assert.Equal(t, []interface {}{4,5,6}, values)
	}
	{
		idList := []int{4,5,6}
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"id": f.OP{
						f.In(idList),
					},
					"age": f.OP{f.Eql(18)},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `age` = ? AND `id` IN (?, ?, ?)", sqlS)
		assert.Equal(t, []interface {}{18,4,5,6}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"name": f.OP{
						f.Custom("LIKE %-?-%", 1),
					},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `name` LIKE %-?-%", sqlS)
		assert.Equal(t, []interface {}{1}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"name": f.OP{
						f.IsNull(),
					},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `name` IS NULL", sqlS)
		var interfaceList []interface {}
		assert.Equal(t, interfaceList, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"name": f.OP{
						f.IsNotNull(),
					},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `name` IS NOT NULL", sqlS)
		var interfaceList []interface {}
		assert.Equal(t, interfaceList, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{
					"time": f.OP{
						f.Day(ge.Time(time.Parse(gtime.LayoutTime, "2018-11-11 00:11:11"))),
					},
				},
			},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `time` >= ? AND `time` <= ?", sqlS)
		assert.Equal(t, []interface {}{"2018-11-11 00:00:00", "2018-11-11 23:59:59"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Group: []string{"name"},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` GROUP BY `name`", sqlS)
		_ = values
		// assert.Equal(t, []interface {}{"2019"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{},
			Group: []string{"name","age"},
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` GROUP BY `name`, `age`", sqlS)
		assert.Equal(t, []interface {}(nil), values)
	}
	{
		qb := f.QB{
			Table: "user",
			Order: []f.Order{ {"name", f.DESC} },
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` ORDER BY `name` DESC", sqlS)
		assert.Equal(t, []interface {}(nil), values)
	}
	{
		qb := f.QB{
			Table: "user",
			Order: []f.Order{ {"name", f.DESC},{"age", f.DESC} },
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` ORDER BY `name` DESC, `age` DESC", sqlS)
		assert.Equal(t, []interface {}(nil), values)
	}
	{
		qb := f.QB{
			Table: "user",
			Order: []f.Order{ {"name", f.DESC},{"age", f.ASC} },
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` ORDER BY `name` DESC, `age` ASC", sqlS)
		assert.Equal(t, []interface {}(nil), values)
	}
	{
		qb := f.QB{
			Table: "user",
			Order: []f.Order{ {"name", f.ASC},{"age", f.DESC} },
		}
		sqlS, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` ORDER BY `name` ASC, `age` DESC", sqlS)
		assert.Equal(t, []interface {}(nil), values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: []f.AND{
				{"id": f.OP{f.Eql("13")},},
			},
			Update: f.Map{
				"name": "sam",
			},
		}
		query, values := qb.GetUpdate()
		assert.Equal(t, "UPDATE `user` SET `name` = ? WHERE `id` = ?", query)
		assert.Equal(t, []interface {}{"sam", "13"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Insert: f.Map{
				"name": "sam",
			},
		}
		insert, values := qb.GetInsert()
		assert.Equal(t, "INSERT INTO `user` (`name`) VALUES (?)", insert)
		assert.Equal(t, []interface {}{"sam"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Insert: f.Map{
				"name": "sam",
				"age": 19,
			},
		}
		insert, values := qb.GetInsert()
		assert.Equal(t, "INSERT INTO `user` (`age`, `name`) VALUES (?, ?)", insert)
		assert.Equal(t, []interface {}{19, "sam"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Count: true,
		}
		insert, values := qb.GetSelect()
		assert.Equal(t, "SELECT count(*) FROM `user`", insert)
		assert.Equal(t, []interface {}(nil), values)
	}
	{
		qb := f.QB{
			Table: "user",
			Order: []f.Order{{"name" , f.ASC}},
			SoftDelete: "deleted_at",
		}
		sql, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `deleted_at` IS NULL ORDER BY `name` ASC", sql)
		assert.Equal(t, []interface {}(nil), values)
	}
	{
		qb := f.QB{
			Table: "user",
			Order: []f.Order{{"name" , f.ASC}},
			SoftDelete: "deleted_at",
			Limit: 1,
		}
		sql, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `deleted_at` IS NULL ORDER BY `name` ASC LIMIT ?", sql)
		assert.Equal(t, []interface {}{1}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Order: []f.Order{{"name" , f.ASC}},
			SoftDelete: "deleted_at",
			Offset: 1,
		}
		sql, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `deleted_at` IS NULL ORDER BY `name` ASC OFFSET ?", sql)
		assert.Equal(t, []interface {}{1}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where:[]f.AND{
				{
					"created_at": f.OP{
						f.CustomSQL("`created_at` < ? OR `created_at` > ?", "2019-11-11", "2019-11-11"),
					},
					"id": f.OP{f.In([]string{"1","2"})},
				},
			},
			SoftDelete: "deleted_at",
		}
		sql, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE (`created_at` < ? OR `created_at` > ?) AND `id` IN (?, ?) AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"2019-11-11", "2019-11-11", "1", "2"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.And(
				"created_at", f.CustomSQL("`created_at` < ? OR `created_at` > ?", "2019-11-11", "2019-11-11"),
				"id", f.In([]string{"1","2"}),
				),
			SoftDelete: "deleted_at",
		}
		sql, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE (`created_at` < ? OR `created_at` > ?) AND `id` IN (?, ?) AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"2019-11-11", "2019-11-11", "1", "2"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.And(
				"created_at", f.CustomSQL("`created_at` < ? OR `end_time` > ?", "2019-11-11", "2019-11-11"),
				"id", f.In([]string{"1","2"}),
			),
			SoftDelete: "deleted_at",
		}
		sql, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE (`created_at` < ? OR `end_time` > ?) AND `id` IN (?, ?) AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"2019-11-11", "2019-11-11", "1", "2"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.And(
				"created_at", f.CustomSQL("`created_at` < ? OR `end_time` > ?", "2019-11-11", "2019-11-11"),
				"other", f.CustomSQL("`other` = ? OR `other` = ?", "nimo", "tim"),
				"id", f.In([]string{"1","2"}),
			),
			SoftDelete: "deleted_at",
		}
		sql, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE (`created_at` < ? OR `end_time` > ?) AND `id` IN (?, ?) AND (`other` = ? OR `other` = ?) AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"2019-11-11", "2019-11-11", "1", "2", "nimo", "tim"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where: f.And(
				"", f.CustomSQL("`start_time` < ? OR `end_time` > ?", "2019-11-11", "2019-11-11"),
				"start_time", f.Eql("2011-01-01"),
				"other", f.CustomSQL("`other` = ? OR `other` = ?", "nimo", "tim"),
				"id", f.In([]string{"1","2"}),
			),
			SoftDelete: "deleted_at",
		}
		sql, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE (`start_time` < ? OR `end_time` > ?) AND `id` IN (?, ?) AND (`other` = ? OR `other` = ?) AND `start_time` = ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"2019-11-11", "2019-11-11", "1", "2", "nimo", "tim", "2011-01-01"}, values)
	}
	{
		qb := f.QB{
			Table: "user",
			Where:[]f.AND{
				{
					"created_at": f.OP{
						f.CustomSQL("`created_at` < ? OR `created_at` > ?", "2019-11-11", "2019-11-11"),
						f.Eql("2010-11-11"),
					},
					"id": f.OP{f.In([]string{"1","2"})},
				},
			},
			SoftDelete: "deleted_at",
		}
		sql, values := qb.GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE (`created_at` < ? OR `created_at` > ?) AND `created_at` = ? AND `id` IN (?, ?) AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"2019-11-11", "2019-11-11", "2010-11-11", "1", "2"}, values)
	}
	{
		sql, values := f.QB{
			Where: f.And("id", 1),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `id` = ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{1}, values)
	}
	{
		{
			sql, values := f.QB{
				Table: "user",
				Where: f.And("title", f.IgnoreEmpty(f.Eql, "")),
			}.GetSelect()
			assert.Equal(t, "SELECT * FROM `user`", sql)
			assert.Equal(t, []interface {}(nil), values)
		}
		{
			sql, values := f.QB{
				Where: f.And("title", f.IgnoreEmpty(f.Eql, "")),
			}.BindModel(&User{}).GetSelect()
			assert.Equal(t, "SELECT * FROM `user` WHERE `deleted_at` IS NULL", sql)
			assert.Equal(t, []interface {}(nil), values)
		}
		{
			sql, values := f.QB{
				Where: f.And("title", "abc"),
			}.BindModel(&User{}).GetSelect()
			assert.Equal(t, "SELECT * FROM `user` WHERE `title` = ? AND `deleted_at` IS NULL", sql)
			assert.Equal(t, []interface{}{"abc"}, values)
		}
	}
	{
		status := "all"
		sql, values := f.QB{
			Where: f.And(
				"status", f.IgnorePattern(f.Eql, status, "all"),
				"gift_id", f.In([]string{"1","2"}),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `gift_id` IN (?, ?) AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"1", "2"}, values)
	}
	{
		status := "done"
		sql, values := f.QB{
			Where: f.And(
				"status", f.IgnorePattern(f.Eql, status, "all"),
				"gift_id", f.In([]string{"1","2"}),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `gift_id` IN (?, ?) AND `status` = ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"1", "2", "done"}, values)
	}
	{
		list := []string{}
		sql, values := f.QB{
			Where: f.And(
				"status", f.Ignore(f.In(list), len(list) == 0),
				"gift_id", f.In([]string{"1","2"}),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `gift_id` IN (?, ?) AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"1", "2"}, values)
	}
	{
		list := []string{"a","b"}
		sql, values := f.QB{
			Where: f.And(
				"status", f.Ignore(f.In(list), len(list) == 0),
				"gift_id", f.In([]string{"1","2"}),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `gift_id` IN (?, ?) AND `status` IN (?, ?) AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"1", "2", "a", "b"}, values)
	}
	{
		sql, values := f.QB{
			Where: f.And(
				"create_at", f.TimeRange(gtime.Range{
					Type:   gtime.Range{}.Type.Enum().Day,
					Start: gtime.ParseChina(gtime.LayoutDate, "2019-11-11"),
					End:  gtime.ParseChina(gtime.LayoutDate, "2019-11-23"),
				}, gtime.LayoutTime),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `create_at` >= ? AND `create_at` <= ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"2019-11-11 00:00:00", "2019-11-23 23:59:59"}, values)
	}
	{
		sql, values := f.QB{
			Where: f.And(
				"date", f.TimeRange(gtime.Range{
					Type:  gtime.Range{}.Type.Enum().Day,
					Start: gtime.ParseChina(gtime.LayoutDate, "2019-11-11"),
					End:  gtime.ParseChina(gtime.LayoutDate, "2019-11-23"),
				}, gtime.LayoutDate),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `date` >= ? AND `date` <= ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"2019-11-11", "2019-11-23"}, values)
	}
	{
		sql, values := f.QB{
			Where: f.And(
				"date", f.TimeRange(gtime.Range{
					Type:   gtime.Range{}.Type.Enum().Month,
					Start: gtime.ParseChina(gtime.LayoutYearAndMonth, "2019-11"),
					End:  gtime.ParseChina(gtime.LayoutYearAndMonth, "2020-02"),
				}, gtime.LayoutDate),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `date` >= ? AND `date` <= ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"2019-11-01", "2020-02-29"}, values)
	}
	{
		sql, values := f.QB{
			Where: f.And(
				"month", f.TimeRange(gtime.Range{
					Type:   gtime.Range{}.Type.Enum().Month,
					Start: gtime.ParseChina(gtime.LayoutYearAndMonth, "2019-11"),
					End:  gtime.ParseChina(gtime.LayoutYearAndMonth, "2020-02"),
				}, gtime.LayoutYearAndMonth),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `month` >= ? AND `month` <= ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"2019-11", "2020-02"}, values)
	}
	{
		sql, values := f.QB{
			Where: f.And(
				"age", f.BetweenInt(18, 65),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `age` BETWEEN ? AND ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{18,65}, values)
	}
	{
		sql, values := f.QB{
			Where: f.And(
				"age", f.BetweenFloat(float64(10)/3, float64(20)/3),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `age` BETWEEN ? AND ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{3.3333333333333335,6.666666666666667}, values)
	}
	{
		sql, values := f.QB{
			Where: f.And(
				"age", f.BetweenString("nimo", "nimoz"),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` WHERE `age` BETWEEN ? AND ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"nimo","nimoz"}, values)
	}
	{
		sql, values := f.QB{
			Table: "goods",
			Where: f.And(
				"saleCount", f.Custom("= `inventory`"),
			),
		}.GetSelect()
		assert.Equal(t, "SELECT * FROM `goods` WHERE `saleCount` = `inventory`", sql)
		assert.Equal(t, []interface{}(nil), values)
	}
	{
		sql, values := f.QB{
			Table: "goods",
			Where: f.And(
				"saleCount", f.IgnoreFilter(),
			),
		}.GetSelect()
		assert.Equal(t, "SELECT * FROM `goods`", sql)
		assert.Equal(t, []interface{}(nil), values)
	}
	{
		sql, values := f.QB{
			Table: "goods",
			Where: f.And(
				"age", f.Eql(1),
				"saleCount", f.IgnoreFilter(),
			),
		}.GetSelect()
		assert.Equal(t, "SELECT * FROM `goods` WHERE `age` = ?", sql)
		assert.Equal(t, []interface{}{1}, values)
	}
	as.PanicError("f.Filter is empty struct", func() {
		sql, values := f.QB{
			Table: "goods",
			Where: f.And(
				"age", f.Filter{},
			),
		}.GetSelect()
		assert.Equal(t, "SELECT * FROM `goods`", sql)
		assert.Equal(t, []interface{}(nil), values)
	})
	{
		sql, values := f.QB{
			UseIndex: "abc",
			Where: f.And(
				"status", f.Eql("all"),
			),
		}.BindModel(&User{}).GetSelect()
		assert.Equal(t, "SELECT * FROM `user` USE INDEX(`abc`) WHERE `status` = ? AND `deleted_at` IS NULL", sql)
		assert.Equal(t, []interface{}{"all"}, values)
	}
}
