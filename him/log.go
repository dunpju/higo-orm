package him

import "fmt"

var (
	enums map[string]*impl
)

const (
	// Silent silent log level
	Silent logLevel = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
)

func init() {
	enums = make(map[string]*impl)
	enums["Silent"] = newEnum(int(Silent), "Silent")
	enums["Error"] = newEnum(int(Error), "Error")
	enums["Warn"] = newEnum(int(Warn), "Warn")
	enums["Info"] = newEnum(int(Info), "Info")
}

func LogLevel(level string) (*impl, error) {
	if e, ok := enums[level]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("%s log level undefined", level)
}

type logLevel int

func (this logLevel) enum() *impl {
	if e, ok := enums[this.Level()]; ok {
		return e
	} else {
		panic(fmt.Errorf("%d log level undefined", this))
	}
}

func (this logLevel) Code() int {
	return this.enum().code
}

func (this logLevel) Level() string {
	return this.enum().level
}

type impl struct {
	code  int
	level string
}

func newEnum(code int, level string) *impl {
	return &impl{code: code, level: level}
}
