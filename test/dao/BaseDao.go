package dao

import (
	"github.com/dunpju/higo-orm/arm"
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-throw/exception"
	"github.com/go-sql-driver/mysql"
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

type Transaction struct {
	tx *gorm.DB
}

func newTransaction(tx *gorm.DB) *Transaction {
	return &Transaction{tx: tx}
}

func (this *Transaction) Transaction(fn func(tx *gorm.DB) error) {
	err := this.tx.Transaction(func(tx *gorm.DB) (err error) {
		err = fn(tx)
		if err == nil {
			tx.Commit()
		}
		return
	})
	if err != nil {
		if e, ok := err.(*mysql.MySQLError); ok {
			exception.Throw(exception.Message(e.Message),
				exception.Code(int(e.Number)),
				exception.Data(nil))
		} else {
			panic(err)
		}
	}
}

func (this *BaseDao) CheckError(gormDB *gorm.DB) {
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
}
