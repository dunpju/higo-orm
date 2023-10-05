package stubs

import (
	"io"
	"os"
	"path"
	"runtime"
)

type Stub struct {
	filename string
	context  string
}

func NewStub(filename string) *Stub {
	return &Stub{filename: filename, context: context(filename)}
}

func (this *Stub) Context() string {
	return this.context
}

func context(filename string) string {
	_, file, _, _ := runtime.Caller(0)
	file = path.Dir(file) + string(os.PathSeparator) + filename
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	cxt, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(cxt)
}
