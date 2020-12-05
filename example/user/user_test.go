package userDataStorage

import (
	"context"
	f "github.com/og/gofree"
	connectRDS "github.com/og/gofree/example/connect/rds"
	m "github.com/og/gofree/example/model"
	"github.com/og/x/test"
	"testing"
)

func TestDataStorage_User(t *testing.T) {
	as := gtest.NewAS(t)
	rds, err := connectRDS.NewRDS() ; as.NoError(err)
	ds := NewDataStorage(rds)
	ctx := context.Background() // 一般是在 http.Request{} 中拿到的 Context()  r.Context()
	id := f.UUID()
	values := []interface{}{id,"13611112222","nimo",18,0,"2020-12-05 15:14:35","2020-12-05 15:14:35", nil}
	_, err = rds.Main.Core.Exec("INSERT INTO `user` (`id`, `mobile`, `name`, `age`, `disabled`, `created_at`, `updated_at`, `deleted_at`) VALUES (?,?,?,?,?,?,?,?)", values...)
	as.NoError(err)
	user, hasUser, err := ds.User(ctx, m.IDUser(id))
	as.NoError(err)
	as.Equal(hasUser, true)
	{
		as.Equal(id, user.ID.String())
		as.Equal(user.Name, "nimo")
		as.Equal(user.Age, 18)
		as.Equal(user.Disabled, false)
	}
	user, err = ds.UserMustHas(ctx, user.ID) ; as.NoError(err)
	{
		as.Equal(id, user.ID.String())
		as.Equal(user.Name, "nimo")
		as.Equal(user.Age, 18)
		as.Equal(user.Disabled, false)
	}
	{
		notExistID := f.UUID()
		_, err = ds.UserMustHas(ctx, m.IDUser(notExistID))
		as.ErrorString(err, "user " + notExistID +" 不存在")
	}
	{
		mobile, hasUser, err := ds.UserMobile(ctx, user.ID) ; as.NoError(err)
		as.Equal(mobile, "13611112222")
		as.Equal(hasUser, true)
	}
	{
		_, hasUser, err := ds.UserMobile(ctx, m.IDUser(f.UUID())) ; as.NoError(err)
		as.Equal(hasUser, false)
	}
}