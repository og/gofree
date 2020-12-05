package userDataStorage

import (
	"context"
	"errors"
	f "github.com/og/gofree"
	m "github.com/og/gofree/example/model"
)

func (dep DataStorage) User(ctx context.Context, id m.IDUser) (user m.User, hasUser bool, reject error) {
	reject = dep.rds.Main.OneID(ctx, &user, &hasUser, id) ; if reject != nil {return}
	return
}

func (dep DataStorage) UserMustHas(ctx context.Context, id m.IDUser) (user m.User, reject error) {
	var hasUser bool
	user, hasUser, reject = dep.User(ctx, id) ; if reject != nil {return}
	if !hasUser {
		reject = errors.New("user " + id.String() + " 不存在") ; return
	}
	return
}
func (dep DataStorage) UserMobile(ctx context.Context, id m.IDUser) (mobile string, hasUser bool, reject error) {
	userCol := m.User{}.Column()
	reject = dep.rds.Main.QueryRowScan(ctx, &hasUser, f.QB{
		Table: m.User{}.TableName(),
		Select: []f.Column{userCol.Mobile},
		Where: f.And(userCol.ID, id),
		Check: []string{"SELECT `mobile` FROM `user` WHERE `id` = ? AND `deleted_at` IS NULL"},
	}, &mobile) ; if reject != nil {return}
	return
}
