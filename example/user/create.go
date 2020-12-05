package userDataStorage

import (
	"context"
	m "github.com/og/gofree/example/model"
	IUserDataStorage "github.com/og/gofree/example/user/interface"
)

func (dep DataStorage) Create(ctx context.Context, data IUserDataStorage.Create) (user m.User, reject error) {
	user = m.User{
		Name: data.Name,
		Mobile: data.Mobile,
		Age: data.Age,
		Disabled: data.Disabled,
	}
	reject = dep.rds.Main.Create(ctx, &user) ; if reject != nil {return}
	return
}