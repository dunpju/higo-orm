package arm

type Flag int

type IEntity interface {
	IsEdit() bool
	Edit(edit bool)
	SetFlag(flag Flag)
	IsFlag(flag Flag) bool
	PrimaryEmpty() bool
}

func (this Flag) Apply(entity IEntity) {
	entity.SetFlag(this)
}
