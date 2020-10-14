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
