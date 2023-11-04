package arm

type Flag int

type IEntity interface {
	IsEdit() bool
	Edit(edit bool)
	Flag(flag Flag)
	Equals(flag Flag) bool
	PrimaryEmpty() bool
}

func (this Flag) Apply(entity IEntity) {
	entity.Flag(this)
}
