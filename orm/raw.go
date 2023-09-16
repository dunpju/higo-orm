package orm

type raw struct {
	sql  string
	args []interface{}
	err  error
}

func (this raw) ToSql() (string, []interface{}, error) {
	return this.sql, this.args, this.err
}

type whereRaw struct {
	sql  string
	args []interface{}
	err  error
}

func (this whereRaw) ToSql() (string, []interface{}, error) {
	return this.sql, this.args, this.err
}
