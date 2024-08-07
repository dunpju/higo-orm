package him

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"gorm.io/gorm"
)

type SelectBuilder struct {
	db          *DB
	connect     *connect
	countColumn []string
	sumColumn   []string
	isRaw       bool
	isWhereRaw  bool
	columns     []string
	from        string
	joins       []join
	wheres      *Wheres
	hasOffset   bool
	offset      uint64
	hasLimit    bool
	limit       uint64
	hasOrderBy  bool
	orderBy     []string
	hasGroupBys bool
	groupBys    []string
	hasHaving   bool
	havings     []having
	hasColumn   bool
	column      []column
	hasDistinct bool
	Error       error
}

func newSelectBuilder(connection string) *SelectBuilder {
	var (
		dbc *connect
		err error
	)

	if connection != "" {
		dbc, err = getConnect(connection)
		if err != nil {
			return &SelectBuilder{Error: err}
		}
	} else {
		dbc, err = getConnect(DefaultConnect)
		if err != nil {
			return &SelectBuilder{Error: err}
		}
	}
	return query(dbc)
}

func query(dbc *connect) *SelectBuilder {
	return &SelectBuilder{
		db:          dbc.db,
		connect:     dbc,
		countColumn: make([]string, 0),
		sumColumn:   make([]string, 0),
		columns:     make([]string, 0),
		joins:       make([]join, 0),
		wheres:      NewWheres(),
		orderBy:     make([]string, 0),
		groupBys:    make([]string, 0),
		havings:     make([]having, 0),
		column:      make([]column, 0),
	}
}

func (this *SelectBuilder) Connect() *connect {
	return this.connect
}

func (this *SelectBuilder) DB() *gorm.DB {
	return this.db.GormDB()
}

func (this *SelectBuilder) clone() *SelectBuilder {
	return &SelectBuilder{
		db:          this.db,
		connect:     this.connect,
		countColumn: this.countColumn,
		sumColumn:   this.sumColumn,
		isRaw:       this.isRaw,
		isWhereRaw:  this.isWhereRaw,
		columns:     this.columns,
		from:        this.from,
		joins:       this.joins,
		wheres:      this.wheres,
		hasOffset:   this.hasOffset,
		offset:      this.offset,
		hasLimit:    this.hasLimit,
		limit:       this.limit,
		hasOrderBy:  this.hasOrderBy,
		orderBy:     this.orderBy,
		hasGroupBys: this.hasGroupBys,
		groupBys:    this.groupBys,
		hasHaving:   this.hasHaving,
		havings:     this.havings,
		hasColumn:   this.hasColumn,
		column:      this.column,
		hasDistinct: this.hasDistinct,
		Error:       this.Error,
	}
}

func (this *SelectBuilder) begin(db *gorm.DB) *SelectBuilder {
	this.db.gormDB = db
	return this
}

func (this *SelectBuilder) _select(columns ...string) *SelectBuilder {
	this.columns = append(this.columns, columns...)
	return this
}

func (this *SelectBuilder) _from(from string) *SelectBuilder {
	this.from = from
	return this
}

func (this *SelectBuilder) Query() *Select {
	conn, err := getConnect(this.connect.DB().connect)
	if err != nil {
		this.Error = err
	} else {
		this.db = newDB(conn.DB().GormDB(), this.connect.db.connect, conn.dbc, this.db.begin)
		return newSelect(this.db, this.db.gormDB)
	}
	return newSelect(this.db, nil)
}

func (this *SelectBuilder) Offset(offset uint64) *SelectBuilder {
	this.hasOffset = true
	this.offset = offset
	return this
}

func (this *SelectBuilder) Limit(limit uint64) *SelectBuilder {
	this.hasLimit = true
	this.limit = limit
	return this
}

func (this *SelectBuilder) OrderBy(orderBys ...any) *SelectBuilder {
	this.hasOrderBy = true
	this.orderBy = columnsToString(orderBys...)
	return this
}

func (this *SelectBuilder) GroupBy(groupBys ...any) *SelectBuilder {
	this.hasGroupBys = true
	this.groupBys = append(this.groupBys, columnsToString(groupBys...)...)
	return this
}

func (this *SelectBuilder) Column(col interface{}, args ...interface{}) *SelectBuilder {
	this.hasColumn = true
	this.column = append(this.column, column{column: col, args: args})
	return this
}

func (this *SelectBuilder) Distinct() *SelectBuilder {
	this.hasDistinct = true
	return this
}

func (this *SelectBuilder) count() *SelectBuilder {
	this.countColumn = append(this.countColumn, fmt.Sprintf("COUNT(*) AS `%s`", _count_))
	return this
}

func (this *SelectBuilder) sum(columns ...string) *SelectBuilder {
	for i, col := range columns {
		if i == 0 && len(columns) == 1 {
			this.sumColumn = append(this.sumColumn, fmt.Sprintf("SUM(%s) AS `%s`", col, _sum_))
		} else {
			this.sumColumn = append(this.sumColumn, fmt.Sprintf("SUM(%s) AS `%s%d`", col, _sum_, i))
		}
	}
	return this
}

func (this *SelectBuilder) ToSql() (string, []interface{}, error) {
	if this.isWhereRaw || this.isRaw {
		return whereRawHandle(*this.wheres)
	}

	isCount := len(this.countColumn) > 0

	if isCount {
		this.columns = make([]string, 0)
		this.columns = append(this.columns, this.countColumn...)
	}
	if len(this.sumColumn) > 0 {
		this.columns = make([]string, 0)
		this.columns = append(this.columns, this.sumColumn...)
	}
	selectBuilder := squirrel.Select(this.columns...)
	selectBuilder = selectBuilder.From(this.from)
	selectBuilder = joins(selectBuilder, this.joins)
	selectBuilder, err := this.whereHandle(selectBuilder, this.wheres)
	if err != nil {
		return "", nil, err
	}
	if this.hasOrderBy && !isCount {
		selectBuilder = selectBuilder.OrderBy(this.orderBy...)
	}
	if this.hasGroupBys {
		selectBuilder = selectBuilder.GroupBy(this.groupBys...)
	}
	if this.hasHaving {
		for _, h := range this.havings {
			selectBuilder = selectBuilder.Having(h.pred, h.rest...)
		}
	}
	if this.hasOffset {
		selectBuilder = selectBuilder.Offset(this.offset)
	}
	if this.hasLimit {
		selectBuilder = selectBuilder.Limit(this.limit)
	}
	if this.hasColumn {
		for _, c := range this.column {
			selectBuilder = selectBuilder.Column(c.column, c.args...)
		}
	}
	if this.hasDistinct {
		selectBuilder = selectBuilder.Distinct()
	}
	return selectBuilder.ToSql()
}
