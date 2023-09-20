package him

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gormDB *gorm.DB
	dbc    *connect
	Error  error
}

func begin(db *DB, tx ...*gorm.DB) *Transaction {
	transaction := &Transaction{}
	dbc, err := getConnect(db.connect)
	if err != nil {
		transaction.Error = err
		return transaction
	}
	transaction.dbc = dbc
	if len(tx) > 0 {
		transaction.gormDB = tx[0]
	} else {
		transaction.gormDB = dbc.DB().GormDB().Begin()
	}

	return transaction
}

func (this *Transaction) Insert() InsertInto {
	return newInsertInto(this.dbc.DB(), this.gormDB)
}

func (this *Transaction) Update() UpdateTable {
	return newUpdateFrom(this.dbc.DB(), this.gormDB)
}

func (this *Transaction) Delete() DeleteFrom {
	return newDeleteFrom(this.dbc.DB(), this.gormDB)
}

func (this *Transaction) Raw(pred string, args ...interface{}) ExecRaw {
	return newExecRaw(this.dbc.DB(), this.gormDB, pred, args)
}

func (this *Transaction) GormDB() *gorm.DB {
	return this.gormDB
}
