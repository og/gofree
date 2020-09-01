package exampleModel

import (
	"github.com/go-sql-driver/mysql"
	f "github.com/og/gofree"
	"time"
)

type IDUser string
type User struct {
	ID IDUser `db:"id"`
	Name string `db:"name"`
	Password string `db:"string"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}
func (user *User) BeforeCreate() {
	user.ID = IDUser(f.UUID())
}
func (user User) TableName() string { return "user" }
func (User) Column() (col struct{
	ID f.Column
	Name f.Column
	Password f.Column
	CreatedAt f.Column
	UpdatedAt f.Column
	DeletedAt f.Column
}) {
	col.ID = "id"
	col.Name = "name"
	col.Password = "password"
	col.CreatedAt = "created_at"
	col.UpdatedAt = "updated_at"
	col.DeletedAt = "deleted_at"
	return col
}

type IDUserWechat string
type UserWechat struct {
	ID IDUserWechat `db:"id"`
	OpenID string `db:"openid"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}
func (user UserWechat) TableName () string {return "user_wechat"}
func (user *UserWechat) BeforeCreate() {
	user.ID = IDUserWechat(f.UUID())
}
func (IDUserWechat) Column() (col struct{
	ID f.Column
	OpenID f.Column
	CreatedAt f.Column
	UpdatedAt f.Column
	DeletedAt f.Column
}) {
	col.ID = "id"
	col.OpenID = "openid"
	col.CreatedAt = "created_at"
	col.UpdatedAt = "updated_at"
	col.DeletedAt = "deleted_at"
	return
}
type UserWithWechat struct {
	UserID IDUser `db:"user.id"`
	Name string `db:"user.name"`
	Password string `db:"user.password"`
	OpenID string `db:"user_wechat.openid"`
}

