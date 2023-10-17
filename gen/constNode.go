package gen

import "go/ast"

type constNodeCollect struct {
	collect []*ast.ValueSpec
}

func newConstNodeCollect() *constNodeCollect {
	return &constNodeCollect{collect: make([]*ast.ValueSpec, 0)}
}

func (this *constNodeCollect) append(valueSpec *ast.ValueSpec) {
	has := false
	var index int
	for i, vs := range this.collect {
		if vs.Names[0].Name == valueSpec.Names[0].Name {
			has = true
			index = i
			break
		}
	}
	if !has {
		this.collect = append(this.collect, valueSpec)
	} else {
		this.collect[index] = valueSpec
	}
}
