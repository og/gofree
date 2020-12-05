package userDataStorage

import (
	connectRDS "github.com/og/gofree/example/connect/rds"
	IUserDataStorage "github.com/og/gofree/example/user/interface"
)

type DataStorage struct {
	rds connectRDS.RDS
}
func NewDataStorage(rds connectRDS.RDS) IUserDataStorage.Interface {
	return DataStorage{
		rds: rds,
	}
}