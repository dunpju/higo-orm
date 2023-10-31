package arm

import (
	"database/sql"
	"github.com/dunpju/higo-orm/him"
)

type BaseDao struct {
	dao IDao
}

func NewBaseDao(dao IDao) *BaseDao {
	return &BaseDao{dao: dao}
}

func (this *BaseDao) Begin(opts ...*sql.TxOptions) *him.TX {
	return this.dao.IModel().Begin(opts...)
}
