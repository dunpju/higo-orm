package him

type IConstructor[T any] interface {
	New(properties ...IProperty) T
}
