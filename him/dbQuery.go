package him

import "gorm.io/gorm"

func (this *DB) Query() Select {
	conn, err := getConnect(this.connect)
	if err != nil {
		this.Error = err
	} else {
		this.db = newDB(conn.DB().GormDB(), this.connect, this.begin)
		return newSelect(this.db, this.gormDB)
	}
	return newSelect(this, nil)
}

type Select struct {
	db         *DB
	gormDB     *gorm.DB
	selectFrom SelectFrom
}

func newSelect(db *DB, gormDB *gorm.DB) Select {
	return Select{db: db, gormDB: gormDB}
}

func (this Select) selectBuilder() SelectBuilder {
	if this.db.begin {
		this.db.Builder = newSelectBuilder(this.db.connect).begin(this.gormDB)
	} else {
		this.db.Builder = newSelectBuilder(this.db.connect)
	}
	return this.db.Builder.(SelectBuilder)
}

func (this Select) Distinct() Select {
	this.db.Builder = this.selectBuilder().Distinct()
	return this
}

func (this Select) Select(columns ...string) SelectFrom {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	this.db.Builder = this.selectBuilder()._select(columns...)
	return newSelectFrom(this.db, this.gormDB)
}

func (this Select) Raw(pred string, args ...interface{}) SelectRaw {
	this.db.Builder = this.selectBuilder().Raw(pred, args...)
	return newSelectRaw(this.db, this.gormDB)
}

type SelectRaw struct {
	db     *DB
	gormDB *gorm.DB
}

func newSelectRaw(db *DB, gormDB *gorm.DB) SelectRaw {
	return SelectRaw{db: db, gormDB: gormDB}
}

func (this SelectRaw) Get(dest interface{}) *gorm.DB {
	return this.db.Builder.(SelectBuilder).Get(dest)
}

type SelectFrom struct {
	db     *DB
	gormDB *gorm.DB
}

func newSelectFrom(db *DB, gormDB *gorm.DB) SelectFrom {
	return SelectFrom{db: db, gormDB: gormDB}
}

func (this SelectFrom) From(from string) SelectBuilder {
	this.db.Builder = this.db.Builder.(SelectBuilder)._from(from)
	return this.db.Builder.(SelectBuilder)
}
