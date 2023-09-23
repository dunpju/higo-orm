package him

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

type BaseDao struct {
	db  *DB
	err error
}

func NewBaseDao(connect string) (*BaseDao, error) {
	return newBaseDao(connect)
}

func newBaseDao(connect string) (*BaseDao, error) {
	conn, err := DBConnect(connect)
	if err != nil {
		return nil, err
	}
	return &BaseDao{db: conn}, nil
}

func (this *BaseDao) DB() *DB {
	return this.db
}

func (this *BaseDao) GormDB() *gorm.DB {
	return this.db.GormDB()
}

func (this *BaseDao) Error() error {
	return this.err
}

func (this *BaseDao) BeginTX(opts ...*sql.TxOptions) *gorm.DB {
	baseDao, err := newBaseDao(this.db.connect)
	if err != nil {
		this.err = err
		return nil
	}
	return baseDao.db.GormDB().Begin(opts...)
}

func (this *BaseDao) Begin(opts ...*sql.TxOptions) *TX {
	return newTX(this.BeginTX(opts...))
}

type TX struct {
	tx *gorm.DB
}

func newTX(tx *gorm.DB) *TX {
	return &TX{tx: tx}
}

func (this *TX) Transaction(fn func(tx *gorm.DB) error) error {
	return this.tx.Transaction(func(tx *gorm.DB) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%s", r)
			}
		}()
		err = fn(tx)
		if err == nil {
			return tx.Commit().Error
		}
		return
	})
}
