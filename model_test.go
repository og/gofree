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
	IsSuper bool `db:"is_super"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
const (
	fUserID f.Field  = "id"
	fUserName f.Field = "name"
	fUserIsSuper f.Field = "is_super"
	fUserCreatedAt f.Field = "created_at"
	fUserUpdatedAt f.Field = "updated_at"
	fUserDeletedAt f.Field = "deleted_at"
)

func (User) TableName() string {
	return "user"
}
func (model *User) BeforeCreate() {
	if model.ID == "" {
		model.ID = IDUser(f.UUID())
	}
}

/*
CREATE TABLE `goods` (
  `id` char(36) NOT NULL DEFAULT '',
  `title` varchar(40) NOT NULL DEFAULT '',
  `sale_quantity` int(11) NOT NULL,
	`inventory` int(11) NOT NULL,
	`price` decimal(11,2) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/
type IDGoods string
type Goods struct {
	ID IDGoods `db:"id"`
	Title string `db:"title"`
	SaleQuantity int `db:"sale_quantity"`
	Inventory int `db:"inventory"`
	Price float64 `db:"price"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
func (Goods) Field () (field struct {
	ID f.Field
	Title f.Field
	SaleQuantity f.Field
	Inventory f.Field
	Price f.Field
	CreatedAt f.Field
	UpdatedAt f.Field
	DeletedAt f.Field
}) {
	field.ID = "id"
	field.Title = "title"
	field.SaleQuantity = "sale_quantity"
	field.Inventory = "inventory"
	field.Price = "price"
	field.CreatedAt = "created_at"
	field.UpdatedAt = "updated_at"
	field.DeletedAt = "deleted_at"
	return
}
func (Goods) TableName() string { return "goods" }
func (model *Goods) BeforeCreate() {
	if model.ID == "" {
		model.ID = IDGoods(f.UUID())
	}
}

/*
CREATE TABLE `goods_detail` (
  `id` char(36) NOT NULL DEFAULT '',
	`goods_id` char(36) NOT NULL DEFAULT '',
  `banner` varchar(4000) NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
  KEY `goods_id` (`goods_detail`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/
type IDGoodsDetail string
type GoodsDetail struct {
	ID      IDGoodsDetail `db:"id"`
	GoodsID IDGoods         `db:"goods_id"`
	Banner  string          `db:"banner"`
}
func (GoodsDetail) Field()(field struct {
	ID f.Field
	GoodsID f.Field
	Banner f.Field
}) {
	field.ID = "id"
	field.GoodsID = "goods_id"
	field.Banner = "banner"
	return
}
func (GoodsDetail) TableName() string { return "goods_detail" }
func (model *GoodsDetail) BeforeCreate() {
	if model.ID == "" {
		model.ID = IDGoodsDetail(f.UUID())
	}
}

type GoodsRelation struct {
	ID IDGoods `db:"goods.id"`
	Title string `db:"goods.title"`
	SaleQuantity string `db:"goods.sale_quantity"`
	Inventory int `db:"goods.inventory"`
	Price string `db:"goods.price"`
	Banner  string `db:"goods_detail.banner"`
}
func (GoodsRelation) Field() (field struct {
	Goods struct {
		ID f.Field
		Title f.Field
		SaleQuantity f.Field
		Inventory f.Field
		Price f.Field
	}
	GoodsDetail struct{
		Banner f.Field
		ID f.Field
		GoodsID f.Field
	}
}) {
	field.Goods.ID = "goods.id"
	field.Goods.Title = "goods.title"
	field.Goods.SaleQuantity = "goods.sale_quantity"
	field.Goods.Inventory = "goods.inventory"
	field.Goods.Price = "goods.price"
	field.GoodsDetail.Banner = "goods_detail.banner"
	field.GoodsDetail.ID = "goods_detail.id"
	field.GoodsDetail.GoodsID = "goods_detail.goods_id"
	return
}
func (goodsRelation GoodsRelation) Relation() (tableName string, join []f.Join) {
	return Goods{}.TableName(), []f.Join{
		f.InnerJoin(
			&Goods{},Goods{}.Field().ID,
			&GoodsDetail{}, GoodsDetail{}.Field().GoodsID,
			),
	}
}