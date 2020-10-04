package exampleGofree

import (
	"github.com/go-sql-driver/mysql"
	f "github.com/og/gofree"
	"time"
)

type IDUser string
func (id IDUser) String() string {return string(id)}
type User struct {
	ID IDUser `db:"id"`
	Name string `db:"name"`
	Age int `db:"age"`
	Disabled bool `db:"disabled"`
	Password string `db:"string"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}
func (user *User) BeforeCreate() {
	if user.ID == "" { user.ID = IDUser(f.UUID()) }
}
func (user User) TableName() string { return "user" }
func (User) Column() (col struct{
	ID f.Column
	Name f.Column
	Age f.Column
	Disabled f.Column
	Password f.Column
	CreatedAt f.Column
	UpdatedAt f.Column
	DeletedAt f.Column
}) {
	col.ID = "id"
	col.Name = "name"
	col.Age = "age"
	col.Disabled = "disabled"
	col.Password = "password"
	col.CreatedAt = "created_at"
	col.UpdatedAt = "updated_at"
	col.DeletedAt = "deleted_at"
	return col
}

type IDUserIntegral string
func (id IDUserIntegral) String() string { return string(id)}
type UserIntegral struct {
	ID IDUserIntegral `db:"id"`
	UserID IDUser `db:"user_id"`
	Integral string `db:"Integral"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}
func (user UserIntegral) TableName () string {return "user_integral"}
func (user *UserIntegral) BeforeCreate() {
	user.ID = IDUserIntegral(f.UUID())
}
func (IDUserIntegral) Column() (col struct{
	ID f.Column
	UserID f.Column
	Integral f.Column
	CreatedAt f.Column
	UpdatedAt f.Column
	DeletedAt f.Column
}) {
	col.ID = "id"
	col.UserID = "user_id"
	col.Integral = "integral"
	col.CreatedAt = "created_at"
	col.UpdatedAt = "updated_at"
	col.DeletedAt = "deleted_at"
	return
}

type UserWithIntegral struct {
	UserID IDUser `db:"user.id"`
	Name string `db:"user.name"`
	Password string `db:"user.password"`
	Integral string `db:"user_integral.integral"`
}

