package arm

import (
	"database/sql"
	"github.com/dunpju/higo-orm/him"
	"gorm.io/gorm"
)

type BaseDao struct {
	db  *him.DB
	err error
}

func NewBaseDao(connect string) (*BaseDao, error) {
	return newBaseDao(connect)
}

func newBaseDao(connect string) (*BaseDao, error) {
	db, err := him.DBConnect(connect)
	if err != nil {
		return nil, err
	}
	return &BaseDao{db: db}, nil
}

func (this *BaseDao) DB() *him.DB {
	return this.db
}

func (this *BaseDao) GormDB() *gorm.DB {
	return this.db.GormDB()
}

func (this *BaseDao) Error() error {
	return this.err
}

func (this *BaseDao) BeginTX(opts ...*sql.TxOptions) *gorm.DB {
	baseDao, err := newBaseDao(this.db.Connect())
	if err != nil {
		this.err = err
		return nil
	}
	return baseDao.db.GormDB().Begin(opts...)
}

func (this *BaseDao) Begin(opts ...*sql.TxOptions) *him.TX {
	return him.NewTX(this.BeginTX(opts...))
}
