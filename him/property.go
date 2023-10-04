package him

type IProperty interface {
	Set(obj any)
}

type Property struct {
	fn func(obj any)
}

func SetProperty(fn func(obj any)) *Property {
	return &Property{fn: fn}
}

func (this *Property) Set(obj any) {
	this.fn(obj)
}

type Properties []IProperty

func (this Properties) Apply(obj any) {
	for _, property := range this {
		property.Set(obj)
	}
}
