package him

import "gorm.io/gorm"

func (this *DB) Update() UpdateTable {
	conn, err := getConnect(this.connect)
	if err != nil {
		this.Error = err
	} else {
		this.db = newDB(conn.DB().GormDB(), this.connect, conn.dbc, this.begin)
		return newUpdateFrom(this.db, this.gormDB)
	}
	return newUpdateFrom(this, nil)
}

type UpdateTable struct {
	db     *DB
	gormDB *gorm.DB
}

func newUpdateFrom(db *DB, gormDB *gorm.DB) UpdateTable {
	return UpdateTable{db: db, gormDB: gormDB}
}

func (this UpdateTable) Table(from string) UpdateBuilder {
	if this.db.begin {
		this.db.Builder = newUpdateBuilder(this.db.connect).begin(this.gormDB).update(from)
	} else {
		this.db.Builder = newUpdateBuilder(this.db.connect).update(from)
	}
	return this.db.Builder.(UpdateBuilder)
}
