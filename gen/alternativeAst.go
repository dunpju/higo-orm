package gen

import "go/ast"

type AlternativeAst struct {
	imports    []string
	constNode  *ast.GenDecl
	structNode *ast.GenDecl
	starExprs  []*ast.StarExpr
	fieldsList []*ast.Field
	funcList   []FnDecl
}

func newAlternativeAst() *AlternativeAst {
	return &AlternativeAst{imports: make([]string, 0), starExprs: make([]*ast.StarExpr, 0), fieldsList: make([]*ast.Field, 0), funcList: make([]FnDecl, 0)}
}
