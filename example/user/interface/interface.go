package IUserDataStorage

import (
	"context"
	m "github.com/og/gofree/example/model"
)

type Interface interface {
	Create(ctx context.Context, data Create) (user m.User, reject error)
	User(ctx context.Context, id m.IDUser) (user m.User, hasUser bool, reject error)
	UserMustHas(ctx context.Context, id m.IDUser) (user m.User, reject error)
	UserMobile(ctx context.Context, id m.IDUser) (mobile string, hasUser bool, reject error)
}
type Create struct {
	Mobile string
	Name string
	Age int
	Disabled bool
}
