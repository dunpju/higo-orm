package him

type IPaginate interface {
	SetTotal(total uint64) IPaginate
	SetPerPage(perPage uint64) IPaginate
	SetCurrentPage(currentPage uint64) IPaginate
	SetLastPage(lastPage uint64) IPaginate
	SetItems(items interface{}) IPaginate
	GetTotal() uint64
	GetPerPage() uint64
	GetCurrentPage() uint64
	GetLastPage() uint64
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

func (this *Paginate) GetTotal() uint64 {
	return this.Total
}
func (this *Paginate) GetPerPage() uint64 {
	return this.PerPage
}
func (this *Paginate) GetCurrentPage() uint64 {
	return this.CurrentPage
}
func (this *Paginate) GetLastPage() uint64 {
	return this.LastPage
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

type IPaginateSum interface {
	SetSum(sum string)
	Sum() string
	Dest() interface{}
	Field() interface{}
}

type PaginateSum struct {
	dest        interface{}
	field       interface{}
	sum         string
	Total       uint64      `json:"total"`
	PerPage     uint64      `json:"per_page"`
	CurrentPage uint64      `json:"current_page"`
	LastPage    uint64      `json:"last_page"`
	Items       interface{} `json:"items"`
}

func NewPaginateSum(dest interface{}, field interface{}) *PaginateSum {
	return &PaginateSum{dest: dest, field: field}
}

func (p *PaginateSum) Paginate(paginate IPaginate) IPaginate {
	p.Total = paginate.GetTotal()
	p.PerPage = paginate.GetPerPage()
	p.CurrentPage = paginate.GetCurrentPage()
	p.LastPage = paginate.GetLastPage()
	p.Items = paginate.GetItems()
	return p
}

func (p *PaginateSum) Dest() interface{} {
	return p.dest
}

func (p *PaginateSum) Field() interface{} {
	return p.field
}

func (p *PaginateSum) Sum() string {
	return p.sum
}

func (p *PaginateSum) SetSum(sum string) {
	p.sum = sum
}

func (p *PaginateSum) SetTotal(total uint64) IPaginate {
	p.Total = total
	return p
}

func (p *PaginateSum) SetPerPage(perPage uint64) IPaginate {
	p.PerPage = perPage
	return p
}

func (p *PaginateSum) SetCurrentPage(currentPage uint64) IPaginate {
	p.CurrentPage = currentPage
	return p
}

func (p *PaginateSum) SetLastPage(lastPage uint64) IPaginate {
	p.LastPage = lastPage
	return p
}

func (p *PaginateSum) SetItems(items interface{}) IPaginate {
	p.Items = items
	return p
}

func (p *PaginateSum) GetTotal() uint64 {
	return p.Total
}

func (p *PaginateSum) GetPerPage() uint64 {
	return p.PerPage
}

func (p *PaginateSum) GetCurrentPage() uint64 {
	return p.CurrentPage
}

func (p *PaginateSum) GetLastPage() uint64 {
	return p.LastPage
}

func (p *PaginateSum) GetItems() interface{} {
	return p.Items
}
