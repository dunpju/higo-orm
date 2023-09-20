package him

import "gorm.io/gorm"

type DB struct {
	gormDB  *gorm.DB
	connect string
	prep    *Preprocessor
}

func NewDB(db *gorm.DB, connect string) *DB {
	return &DB{gormDB: db, connect: connect}
}

func (this *DB) Query() SelectBuilder {
	return Query(this.connect)
}

func (this *DB) Insert(into string) InsertBuilder {
	return NewInsertBuilder(this.connect).Insert(into)
}

func (this *DB) Update(table string) UpdateBuilder {
	return NewUpdateBuilder(this.connect).Update(table)
}

func (this *DB) Delete(from string) DeleteBuilder {
	return NewDeleteBuilder(this.connect).Delete(from)
}

func (this *DB) Begin(tx ...*gorm.DB) *Transaction {
	return begin(this.connect, tx...)
}

func (this *DB) TX(tx ...*gorm.DB) *Transaction {
	return this.Begin(tx...)
}

func (this *DB) GormDB() *gorm.DB {
	return this.gormDB
}

func (this *DB) Prep() *Preprocessor {
	return newPreprocessor(this)
}
