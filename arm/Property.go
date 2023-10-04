package arm

import "github.com/dunpju/higo-orm/him"

type Property struct {
	model IModel
}

func newProperty(model IModel) *Property {
	return &Property{model: model}
}

func (this *Property) Property(properties ...him.IProperty) {
	him.Properties(properties).Apply(this.model)
}
