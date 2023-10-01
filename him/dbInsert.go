package him

import (
	"gorm.io/gorm"
)

func (this *DB) Insert() InsertInto {
	conn, err := getConnect(this.connect)
	if err != nil {
		this.Error = err
	} else {
		this.db = newDB(conn.DB().GormDB(), this.connect, this.begin)
		return newInsertInto(this.db, this.gormDB)
	}
	return newInsertInto(this, nil)
}

type InsertInto struct {
	db     *DB
	gormDB *gorm.DB
}

func newInsertInto(db *DB, gormDB *gorm.DB) InsertInto {
	return InsertInto{db: db, gormDB: gormDB}
}

func (this InsertInto) Into(from string) InsertBuilder {
	if this.db.begin {
		this.db.Builder = newInsertBuilder(this.db).begin(this.gormDB).insert(from)
	} else {
		this.db.Builder = newInsertBuilder(this.db).insert(from)
	}
	return this.db.Builder.(InsertBuilder)
}
