package f_test

import (
	"database/sql"
	f "github.com/og/gofree"
	"testing"
	"time"
)
type IDUser string
type User struct {
	ID IDUser `db:"id"`
	Name string `db:"name"`
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
func TestDocOneQB(t *testing.T) {
	db := NewDB()
	query := struct {
		Name string
		Age int
		Gender string
	}{
		Name: "nimo",
		Age: 27,
		Gender: "f",
	}
	// 按50%模拟空 gender
	if time.Now().Second() % 2 == 0 {
		query.Gender = ""
	}
	var foundUser bool
	user := User{}
	db.OneQB(&user, &foundUser, f.QB{
		Where:
		f.And("name", query.Name).
			And("age", query.Age).
			And("gender", f.EqualIgnoreEmpty(query.Gender)),
	}).
	Check("SELECT `id`, `name`, `is_super`, `created_at`, `updated_at`, `deleted_at` FROM `user` WHERE `age` = ? AND `gender` = ? AND `name` = ? AND `deleted_at` IS NULL",
		  "SELECT `id`, `name`, `is_super`, `created_at`, `updated_at`, `deleted_at` FROM `user` WHERE `age` = ? AND `name` = ? AND `deleted_at` IS NULL")
}
