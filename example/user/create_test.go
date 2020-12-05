package userDataStorage

import (
	"context"
	connectRDS "github.com/og/gofree/example/connect/rds"
	IUserDataStorage "github.com/og/gofree/example/user/interface"
	"github.com/og/x/test"
	"testing"
)

func TestDataStorage_Create(t *testing.T) {
		as := gtest.NewAS(t)
		rds, err := connectRDS.NewRDS() ; as.NoError(err)
		ds := NewDataStorage(rds)
		ctx := context.Background() // 一般是在 http.Request{} 中拿到的 Context()  r.Context()
		user, err := ds.Create(ctx, IUserDataStorage.Create{
			Mobile: "13611112222",
			Name: "nimo",
			Age: 18,
			Disabled: false,
		})
		{
			as.NoError(err)
			as.Len(user.ID.String(), 36)
			as.Equal(user.Name, "nimo")
			as.Equal(user.Age, 18)
			as.Equal(user.Disabled, false)
		}
		{
			readUser, err := ds.UserMustHas(ctx, user.ID)
			as.NoError(err)
			as.Equal(readUser.ID, user.ID)
			as.Equal(readUser.Name, "nimo")
			as.Equal(readUser.Age, 18)
			as.Equal(readUser.Disabled, false)
		}
}