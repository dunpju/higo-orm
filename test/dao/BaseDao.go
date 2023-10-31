package dao

import (
	"github.com/dunpju/higo-orm/arm"
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

type BaseDao struct {
	*arm.BaseDao
}

func newBaseDao() *BaseDao {
	dao, err := arm.NewBaseDao(him.DefaultConnect)
	if err != nil {
		panic(err)
	}
	return &BaseDao{dao}
}

func (this *BaseDao) CheckError(gormDB *gorm.DB) {
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
}
