package dao

import (
	"github.com/dunpju/higo-orm/arm"
	"gorm.io/gorm"
)

type BaseDao struct {
	*arm.BaseDao
}

func newBaseDao() *BaseDao {
	return &BaseDao{}
}

func (this *BaseDao) CheckError(gormDB *gorm.DB) {
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
}
