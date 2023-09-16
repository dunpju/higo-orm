package orm

type Paginate struct {
	Total       uint64
	PerPage     uint64
	CurrentPage uint64
	LastPage    uint64
	Items       interface{}
}
