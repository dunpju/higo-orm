package him

type CaseWhen interface {
	Field() string
	Case() string
	WhenThen() []*WhenThen
	ELSE() *Else
}

type WhenThen struct {
	when, then string
}

func NewWhenThen(when, then any) *WhenThen {
	return &WhenThen{when: toString(when), then: toString(then)}
}

type Else struct {
	value string
}

func NewElse(value string) *Else {
	return &Else{value: value}
}
