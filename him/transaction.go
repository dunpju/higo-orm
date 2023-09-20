package him

import (
	"gorm.io/gorm"
)

type Transaction struct {
	db      *gorm.DB
	connect string
	Error   error
}

func begin(connect string, tx ...*gorm.DB) *Transaction {
	transaction := &Transaction{connect: connect}
	if len(tx) > 0 {
		transaction.db = tx[0]
	} else {
		dbc, err := getConnect(connect)
		if err != nil {
			transaction.Error = err
			return transaction
		}
		transaction.db = dbc.DB().GormDB().Begin()
	}
	return transaction
}

func (this *Transaction) Insert(into string) InsertBuilder {
	return newInsertBuilder(this.connect).Transaction(this.db).Insert(into)
}

func (this *Transaction) Update(table ...string) UpdateBuilder {
	return newUpdateBuilder(this.connect).Transaction(this.db).Update(table...)
}

func (this *Transaction) Delete(from ...string) DeleteBuilder {
	return newDeleteBuilder(this.connect).Transaction(this.db).Delete(from...)
}
