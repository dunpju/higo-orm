package him

type prepType int

const (
	PrepInsert prepType = iota + 1
	PrepUpdate
)

type Preprocessor struct {
	db       *DB
	prepType prepType
}

func newPreprocessor(db *DB) *Preprocessor {
	p := &Preprocessor{db: db}
	return p
}

func (this *Preprocessor) Insert(into string) InsertBuilder {
	this.prepType = PrepInsert
	this.db.prep = this
	return this.db.Insert(into)
}

func (this *Preprocessor) Update(table string) UpdateBuilder {
	this.prepType = PrepUpdate
	this.db.prep = this
	return this.db.Update(table)
}
func (this *Preprocessor) setPrepType(pt prepType) *Preprocessor {
	this.prepType = pt
	return this
}

f
