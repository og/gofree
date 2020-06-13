package f_test

import (
	f "github.com/og/gofree"
	"testing"
)

func TestBasic(t *testing.T) {
	db := f.NewDatabase(f.DataSourceName{
		DriverName: "mysql",
		User: "root",
		Password: "somepass",
		Host: "localhost",
		Port: "3306",
		DB: "test_gofree",
	})
	mock := struct {
		toothbrushID IDGoods
	} {
		toothbrushID: IDGoods(f.UUID()),
	}
	{
		f.ResetAndMock(db, f.Mock{
			Tables: []interface{}{
				[]Goods{
					{
						ID: mock.toothbrushID,
						Title:"Bamboo charcoal toothbrush",
						SaleQuantity: 10,
						Inventory: 22,
						Price:1,
					},
				},
				[]GoodsDetail{
					{
						GoodsID: mock.toothbrushID,
						Banner: "https://picsum.photos/333/333",
					},
				},
			},
		})
	}
	goodsRelation := GoodsRelation{}
	has := false

	db.OneRelationQB(&goodsRelation,&has, f.QB{
		Where: f.And(goodsRelation.Field().Goods.ID, mock.toothbrushID),
	})

}
