package gen

import (
	"bytes"
	"go/ast"
)

type funcListCollect struct {
	collect []FnDecl
}

func newFuncListCollect() *funcListCollect {
	return &funcListCollect{collect: make([]FnDecl, 0)}
}

func (this *funcListCollect) append(fd FnDecl) {
	has := false
	var index int
	for i, s := range this.collect {
		if s.Name == fd.Name {
			has = true
			index = i
			break
		}
	}
	if !has {
		this.collect = append(this.collect, fd)
	} else {
		this.collect[index] = fd
	}
}

type FnDecl struct {
	Name string
	Fd   *ast.FuncDecl
}

func newFnDecl(name string, fd *ast.FuncDecl) FnDecl {
	return FnDecl{Name: name, Fd: fd}
}

type FuncDeclWrite struct {
	buf *bytes.Buffer
}

func newFuncDeclWrite() *FuncDeclWrite {
	return &FuncDeclWrite{buf: bytes.NewBufferString("")}
}

func (this *FuncDeclWrite) Write(p []byte) (n int, err error) {
	return this.buf.Write(p)
}
