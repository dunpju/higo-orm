package Transaction

import (
	"github.com/dunpju/higo-orm/orm"
	"github.com/dunpju/higo-orm/orm/Transaction/Insert"
	"github.com/dunpju/higo-orm/orm/Transaction/Update"
	"gorm.io/gorm"
)

type Transaction struct {
	db    *gorm.DB
	Error error
}

func Begin(db ...*gorm.DB) *Transaction {
	tx := &Transaction{}
	if len(db) > 0 {
		tx.db = db[0]
	} else {
		tx.db, tx.Error = orm.Gorm()
	}
	return tx
}

func (this *Transaction) Insert(into string) orm.InsertBuilder {
	return Insert.Transaction(this.db).Into(into)
}

func (this *Transaction) Update(table ...string) orm.UpdateBuilder {
	return Update.Transaction(this.db).Update(table...)
}

func (this *Transaction) Delete(into string) {
}
