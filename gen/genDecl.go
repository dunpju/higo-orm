package gen

import "bytes"

type GenDeclWrite struct {
	buf *bytes.Buffer
}

func newGenDeclWrite() *GenDeclWrite {
	return &GenDeclWrite{buf: bytes.NewBufferString("")}
}

func (this *GenDeclWrite) Write(p []byte) (n int, err error) {
	return this.buf.Write(p)
}
