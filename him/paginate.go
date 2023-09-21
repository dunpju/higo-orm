package him

type Paginate struct {
	Total       uint64      `json:"total"`
	PerPage     uint64      `json:"per_page"`
	CurrentPage uint64      `json:"current_page"`
	LastPage    uint64      `json:"last_page"`
	Items       interface{} `json:"items"`
}
