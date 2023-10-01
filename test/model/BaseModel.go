package model

import "github.com/dunpju/higo-orm/him"

type BaseModel struct {
}

func (this *BaseModel) DB() *him.DB {
	db, err := him.DBConnect(him.DefaultConnect)
	if err != nil {
		panic(err)
	}
	return db
}
