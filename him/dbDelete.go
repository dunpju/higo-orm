package him

import "gorm.io/gorm"

func (this *DB) Delete() DeleteFrom {
	conn, err := getConnect(this.connect)
	if err != nil {
		this.Error = err
	} else {
		this.slaveDB = newDB(conn.DB().GormDB(), this.connect)
		return newDeleteFrom(this.slaveDB, this.gormDB)
	}
	return newDeleteFrom(this, nil)
}

type DeleteFrom struct {
	db     *DB
	gormDB *gorm.DB
}

func newDeleteFrom(db *DB, gormDB *gorm.DB) DeleteFrom {
	return DeleteFrom{db: db, gormDB: gormDB}
}

func (this DeleteFrom) From(from string) DeleteBuilder {
	if this.db.begin {
		this.db.Builder = newDeleteBuilder(this.db.connect).begin(this.gormDB).delete(from)
	} else {
		this.db.Builder = newDeleteBuilder(this.db.connect).delete(from)
	}
	return this.db.Builder.(DeleteBuilder)
}
