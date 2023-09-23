package him

import (
	"gorm.io/gorm"
)

func (this *DB) Insert() InsertInto {
	conn, err := getConnect(this.connect)
	if err != nil {
		this.Error = err
	} else {
		this.slaveDB = newDB(conn.DB().GormDB(), this.connect)
		return newInsertInto(this.slaveDB, this.gormDB)
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
		return newInsertBuilder(this.db).begin(this.gormDB).insert(from)
	} else {
		return newInsertBuilder(this.db).insert(from)
	}
}
