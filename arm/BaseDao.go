package arm

import (
	"database/sql"
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

type BaseDao struct {
	dao IDao
}

func NewBaseDao(dao IDao) *BaseDao {
	return &BaseDao{dao: dao}
}

func (this *BaseDao) Begin(opts ...*sql.TxOptions) *him.TX {
	return this.dao.GetModel().Begin(opts...)
}

func (this *BaseDao) CheckError(gormDB *gorm.DB) {
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
}
