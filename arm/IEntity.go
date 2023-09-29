package arm

type Flag int

type IEntity interface {
	IsEdit() bool
	Edit(edit bool)
	SetFlag(flag Flag)
	Flag() Flag
	PrimaryEmpty() bool
}

func (this Flag) Apply(entity IEntity) {
	entity.SetFlag(this)
}
