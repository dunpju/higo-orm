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

func Begin(tx ...*gorm.DB) *Transaction {
	transaction := &Transaction{}
	if len(tx) > 0 {
		transaction.db = tx[0]
	} else {
		g, err := orm.Gorm()
		if err != nil {
			transaction.Error = err
		}
		transaction.db = g.Begin()
	}
	return transaction
}

func (this *Transaction) Insert(into string) orm.InsertBuilder {
	return Insert.Transaction(this.db).Into(into)
}

func (this *Transaction) Update(table ...string) orm.UpdateBuilder {
	return Update.Transaction(this.db).Update(table...)
}

func (this *Transaction) Delete(into string) {
}
