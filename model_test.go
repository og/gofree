package f_test

import (
	"database/sql"
	f "github.com/og/gofree"
	"time"
)


type IDUser string
type User struct {
	ID IDUser `db:"id"`
	Name string `db:"name"`
	Age int `db:"age"`
	IsSuper bool `db:"is_super"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	NotSQLField string
	EmptyDBTag string `db:""`
}
func (User) TableName() string {
	return "user"
}
func (model *User) BeforeCreate() {
	if model.ID == "" {
		model.ID = IDUser(f.UUID())
	}
}
func (User) Column() (col struct {
	ID f.Column
	Name f.Column
	IsSuper f.Column
	CreatedAt f.Column
	UpdatedAt f.Column
	DeletedAt f.Column
}) {
	col.ID = "id"
	col.Name = "name"
	col.IsSuper = "is_super"
	col.CreatedAt = "created_at"
	col.UpdatedAt = "updated_at"
	col.DeletedAt = "deleted_at"
	return
}

type UserAddress struct {
	UserID IDUser `db:"user_id"`
	Address string `db:"address"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
func (UserAddress) TableName() string {
	return "user_address"
}
func (model *UserAddress) BeforeCreate() {}
func (UserAddress) Column() (col struct {
	UserID f.Column
	Address f.Column
	CreatedAt f.Column
	UpdatedAt f.Column
	DeletedAt f.Column
}) {
	col.UserID = "user_id"
	col.Address = "address"
	col.CreatedAt = "created_at"
	col.UpdatedAt = "updated_at"
	col.DeletedAt = "deleted_at"
	return
}
type UserWithAddress struct {
	UserID IDUser `db:"user.id"`
	Name string `db:"user.name"`
	Age int `db:"user.age"`
	IsSuper bool `db:"user.is_super"`
	Address string `db:"user_address.address"`
}

func (UserWithAddress) TableName() string {return "user"}
func (UserWithAddress) Column() (col struct{
	UserID f.Column
	Name f.Column
	Age f.Column
	IsSuper f.Column
	Address f.Column
}) {
	col.UserID = "user.id"
	col.Name = "user.name"
	col.IsSuper = "user.is_super"
	col.Age = "user.age"
	col.Address = "user_address.address"
	return
}
func (u UserWithAddress) SQLJoin() []f.Join {
	return []f.Join{
		{
			Type: f.LeftJoin,
			TableName: UserAddress{}.TableName(),
			On: []f.Column{"user.id", "user_address.user_id"},
		},
	}
}

type Log struct {
	ID uint `db:"id" dbAutoIncrement:"true"`
	Message string `db:"message"`
}
func (Log) TableName() string {
	return "log"
}
func (model Log) BeforeCreate() {

}
type Log2 struct {
	ID int `db:"id"  dbAutoIncrement:"true"`
	Message string `db:"message"`
}
func (Log2) TableName() string {
	return "log"
}
func (model Log2) BeforeCreate() {

}
type Log3 struct {
	ID string `db:"id"  dbAutoIncrement:"true"`
	Message string `db:"message"`
}
func (Log3) TableName() string {
	return "log"
}
func (model Log3) BeforeCreate() {

}
type Log4 struct {
	ID bool `db:"id"  dbAutoIncrement:"true"`
	Message string `db:"message"`
}
func (Log4) TableName() string {
	return "log"
}
func (model Log4) BeforeCreate() {

}
type Log5 struct {
	ID int `db:"id" dbAutoIncrement:"t"`
	Message string `db:"message"`
}
func (Log5) TableName() string {
	return "log"
}
func (model Log5) BeforeCreate() {

}
type Log6 struct {
	ID int `db:"id" dbAutoIncrement:"false"`
	Message string `db:"message"`
}
func (Log6) TableName() string {
	return "log"
}
func (model Log6) BeforeCreate() {

}
