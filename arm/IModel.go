package arm

import "fmt"

type Property func(class IModel)
type Properties []Property

func (this Properties) Apply(model IModel) {
	for _, property := range this {
		property(model)
	}
}

type IModel interface {
	New() IModel
	Mutate(attrs ...Property) IModel
	Exist() bool
	TableName() TableName
}

type TableName string

func (this TableName) Alias(alias string) string {
	return fmt.Sprintf("%s AS %s", this, alias)
}
