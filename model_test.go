package f_test

import (
	"database/sql"
	f "github.com/og/gofree"
	"time"
)

func NewIDUser(id string) IDUser {
	return IDUser(id)
}
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
type IDUserLocation string
type UserLocation struct {
	ID IDUserLocation `db:"id"`
	Point f.Point `db:"point"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
func (UserLocation) TableName() string {
	return "user_location"
}
func (model *UserLocation) BeforeCreate() {
	if model.ID == "" {
		model.ID = IDUserLocation(f.UUID())
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

