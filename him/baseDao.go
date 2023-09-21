package him

import (
	"database/sql"
	"gorm.io/gorm"
)

type BaseDao struct {
	db *DB
}

func NewBaseDao(connect string) *BaseDao {
	return newBaseDao(connect)
}

func newBaseDao(connect string) *BaseDao {
	conn, err := DBConnect(connect)
	if err != nil {
		panic(err)
	}
	return &BaseDao{db: conn}
}

func (this *BaseDao) DB() *DB {
	return this.db
}

func (this *BaseDao) GormDB() *gorm.DB {
	return this.db.GormDB()
}

func (this *BaseDao) BeginTX(opts ...*sql.TxOptions) *gorm.DB {
	return newBaseDao(this.db.connect).db.GormDB().Begin(opts...)
}

func (this *BaseDao) Begin(opts ...*sql.TxOptions) *TX {
	return newTX(this.BeginTX(opts...))
}

func (this *BaseDao) CheckError(gormDB *gorm.DB) error {
	if gormDB.Error != nil {
		return gormDB.Error
	}
	return nil
}

type TransactionHandle func(tx *gorm.DB) error

type Transactionable interface {
	Transaction() TransactionHandle
}

type TX struct {
	tx *gorm.DB
}

func newTX(tx *gorm.DB) *TX {
	return &TX{tx: tx}
}

func (this *TX) Transaction(fn func(tx *gorm.DB) error) error {
	return this.tx.Transaction(func(tx *gorm.DB) (err error) {
		err = fn(tx)
		if err == nil {
			tx.Commit()
		}
		return
	})
}
