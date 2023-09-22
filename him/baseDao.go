package him

import (
	"database/sql"
	"gorm.io/gorm"
)

type BaseDao struct {
	db *DB
	tx TXHandle
}

func NewBaseDao(connect string, tx TXHandle) *BaseDao {
	return newBaseDao(connect).setTx(tx)
}

func newBaseDao(connect string) *BaseDao {
	conn, err := DBConnect(connect)
	if err != nil {
		panic(err)
	}
	return &BaseDao{db: conn}
}

func (this *BaseDao) setTx(tx TXHandle) *BaseDao {
	this.tx = tx
	return this
}

func (this *BaseDao) DB() *DB {
	return this.db
}

func (this *BaseDao) GormDB() *gorm.DB {
	return this.db.GormDB()
}

func (this *BaseDao) BeginTX(opts ...*sql.TxOptions) *gorm.DB {
	return newBaseDao(this.db.connect).setTx(this.tx).db.GormDB().Begin(opts...)
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

type TXHandle func(tx *gorm.DB) error

type Transactionable interface {
	Transaction() TXHandle
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
			return tx.Commit().Error
		}
		return
	})
}
