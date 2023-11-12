package him

type IPaginate interface {
	SetTotal(total uint64) IPaginate
	SetPerPage(perPage uint64) IPaginate
	SetCurrentPage(currentPage uint64) IPaginate
	SetLastPage(lastPage uint64) IPaginate
	SetItems(items interface{}) IPaginate
	GetItems() interface{}
}

type Paginate struct {
	Total       uint64      `json:"total"`
	PerPage     uint64      `json:"per_page"`
	CurrentPage uint64      `json:"current_page"`
	LastPage    uint64      `json:"last_page"`
	Items       interface{} `json:"items"`
}

func NewPaginate(properties ...IProperty) *Paginate {
	return (&Paginate{}).Property(properties...)
}

func (this *Paginate) Property(properties ...IProperty) *Paginate {
	Properties(properties).Apply(this)
	return this
}

func (this *Paginate) SetTotal(total uint64) IPaginate {
	this.Total = total
	return this
}

func (this *Paginate) SetPerPage(perPage uint64) IPaginate {
	this.PerPage = perPage
	return this
}

func (this *Paginate) SetCurrentPage(currentPage uint64) IPaginate {
	this.CurrentPage = currentPage
	return this
}

func (this *Paginate) SetLastPage(lastPage uint64) IPaginate {
	this.LastPage = lastPage
	return this
}

func (this *Paginate) SetItems(items interface{}) IPaginate {
	this.Items = items
	return this
}

func (this *Paginate) GetItems() interface{} {
	return this.Items
}

func WithTotal(total uint64) IProperty {
	return SetProperty(func(obj any) {
		obj.(*Paginate).Total = total
	})
}

func WithPerPage(perPage uint64) IProperty {
	return SetProperty(func(obj any) {
		obj.(*Paginate).PerPage = perPage
	})
}

func WithCurrentPage(currentPage uint64) IProperty {
	return SetProperty(func(obj any) {
		obj.(*Paginate).CurrentPage = currentPage
	})
}

func WithLastPage(lastPage uint64) IProperty {
	return SetProperty(func(obj any) {
		obj.(*Paginate).LastPage = lastPage
	})
}

func WithItems(items interface{}) IProperty {
	return SetProperty(func(obj any) {
		obj.(*Paginate).Items = items
	})
}
