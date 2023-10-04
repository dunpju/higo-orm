package arm

import "github.com/dunpju/higo-orm/him"

type Property struct {
	model IModel
	err   error
}

func newProperty(model IModel, err error) *Property {
	return &Property{model: model, err: err}
}

func (this *Property) Property(properties ...him.IProperty) error {
	him.Properties(properties).Apply(this.model)
	return this.err
}
