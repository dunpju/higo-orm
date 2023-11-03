package gen

import (
	"bytes"
	"go/ast"
	"go/token"
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
	Name    string
	FileSet *token.FileSet
	Fd      *ast.FuncDecl
}

func newFnDecl(name string, fileSet *token.FileSet, fd *ast.FuncDecl) FnDecl {
	return FnDecl{Name: name, FileSet: fileSet, Fd: fd}
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
