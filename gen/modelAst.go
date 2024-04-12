package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
)

// astForeach 遍历
func (this *Model) newAstEach() *AlternativeAst {
	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, this.outfile, this.stubContext, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	alternativeAst := newAlternativeAst()
	ast.Inspect(astFile, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.GenDecl:
			if n.Tok.IsKeyword() && n.Tok.String() == token.IMPORT.String() {
				for _, spec := range n.Specs {
					alternativeAst.imports = append(alternativeAst.imports, spec.(*ast.ImportSpec).Path.Value)
				}
			} else if n.Tok.IsKeyword() && n.Tok.String() == token.CONST.String() {
				alternativeAst.constNode = n
			} else if n.Specs != nil && len(n.Specs) > 0 {
				if typeSpec, ok := n.Specs[0].(*ast.TypeSpec); ok {
					structType, ok := typeSpec.Type.(*ast.StructType)
					if ok && typeSpec.Name.Obj.Kind.String() == token.TYPE.String() && typeSpec.Name.String() == modelStructName {
						alternativeAst.structNode = n //找到struct node
						fieldsList := structType.Fields.List
						if fieldsList != nil && len(fieldsList) > 0 {
							for _, field := range fieldsList {
								starExpr, ok := field.Type.(*ast.StarExpr)
								if ok && len(field.Names) == 0 { //找到 StarExpr
									alternativeAst.starExprs = append(alternativeAst.starExprs, starExpr)
								} else if len(field.Names) > 0 && !ok {
									alternativeAst.fieldsList = append(alternativeAst.fieldsList, field)
								}
							}
						}
					}
				}
			}
		case *ast.FuncDecl:
			if len(n.Body.List) > 0 {
				for _, stmt := range n.Body.List {
					if returnStmt, ok := stmt.(*ast.ReturnStmt); ok {
						for _, result := range returnStmt.Results {
							if callExpr, ok := result.(*ast.CallExpr); ok {
								for _, arg := range callExpr.Args {
									if funcLit, ok := arg.(*ast.FuncLit); ok {
										for _, s := range funcLit.Body.List {
											if assignStmt, ok := s.(*ast.AssignStmt); ok {
												for _, lh := range assignStmt.Lhs {
													if selectorExpr, ok := lh.(*ast.SelectorExpr); ok {
														if typeAssertExpr, ok := selectorExpr.X.(*ast.TypeAssertExpr); ok {
															if starExpr, ok := typeAssertExpr.Type.(*ast.StarExpr); ok {
																if ident, ok := starExpr.X.(*ast.Ident); ok {
																	if ident.Name == modelStructName && findProperty(selectorExpr.Sel.Name, this.upperProperties) {
																		alternativeAst.funcList = append(alternativeAst.funcList, newFnDecl(selectorExpr.Sel.Name, fileSet, n))
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	//ast.Print(fileSet, astFile)
	return alternativeAst
}

func (this *Model) oldAstEach(alternativeAst *AlternativeAst) {
	oldFd, err := os.OpenFile(this.outfile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	oldContext, err := io.ReadAll(oldFd)
	if err != nil {
		panic(err)
	}
	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, this.outfile, oldContext, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	funcList := newFuncListCollect()
	for _, fd := range alternativeAst.funcList {
		funcList.append(fd)
	}
	hasStarExprArmModel := false
	ast.Inspect(astFile, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.GenDecl:
			if n.Specs != nil && len(n.Specs) > 0 {
				if typeSpec, ok := n.Specs[0].(*ast.TypeSpec); ok {
					structType, structTypeOk := typeSpec.Type.(*ast.StructType)
					if structTypeOk && typeSpec.Name.Obj.Kind.String() == token.TYPE.String() {
						for _, field := range structType.Fields.List {
							if starExpr, starExprOk := field.Type.(*ast.StarExpr); starExprOk {
								after := fmt.Sprintf("%s.%s", starExpr.X.(*ast.SelectorExpr).X.(*ast.Ident).String(),
									starExpr.X.(*ast.SelectorExpr).Sel.String())
								if starExprArmModel == after {
									hasStarExprArmModel = true // arm.Model
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	newFileBuf := bytes.NewBufferString("")
	inspect := NewInspect(alternativeAst)
	ast.Inspect(astFile, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.File:
			newFileBuf.WriteString(fmt.Sprintf("package %s\n", n.Name.Name))
			newFileBuf.WriteString("\n")
		case *ast.GenDecl:
			inspect.Import(newFileBuf, n)
			inspect.Const(newFileBuf, n)
			inspect.Type(newFileBuf, n, hasStarExprArmModel)
		case *ast.FuncDecl:
			isWithFunc := false
			if len(n.Body.List) > 0 {
				for _, stmt := range n.Body.List {
					if returnStmt, ok := stmt.(*ast.ReturnStmt); ok {
						for _, result := range returnStmt.Results {
							if callExpr, ok := result.(*ast.CallExpr); ok {
								for _, arg := range callExpr.Args {
									if funcLit, ok := arg.(*ast.FuncLit); ok {
										for _, s := range funcLit.Body.List {
											if assignStmt, ok := s.(*ast.AssignStmt); ok {
												for _, lh := range assignStmt.Lhs {
													if selectorExpr, ok := lh.(*ast.SelectorExpr); ok {
														if typeAssertExpr, ok := selectorExpr.X.(*ast.TypeAssertExpr); ok {
															if starExpr, ok := typeAssertExpr.Type.(*ast.StarExpr); ok {
																if ident, ok := starExpr.X.(*ast.Ident); ok {
																	if ident.Name == modelStructName && findProperty(selectorExpr.Sel.Name, this.upperProperties) {
																		funcList.append(newFnDecl(selectorExpr.Sel.Name, fileSet, n))
																		isWithFunc = true
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			if !isWithFunc {
				funcDeclWrite := newFuncDeclWrite()
				err = printer.Fprint(funcDeclWrite, fileSet, n)
				if err != nil {
					panic(err)
				}
				funcDeclWrite.buf.WriteString("\n")
				newFileBuf.WriteString(funcDeclWrite.buf.String())
				newFileBuf.WriteString("\n")
			}
		}
		return true
	})

	for _, fd := range funcList.collect {
		funcDeclWrite := newFuncDeclWrite()
		err = printer.Fprint(funcDeclWrite, fd.FileSet, fd.Fd)
		if err != nil {
			panic(err)
		}
		funcDeclWrite.buf.WriteString("\n")
		newFileBuf.WriteString(funcDeclWrite.buf.String())
		newFileBuf.WriteString("\n")
	}
	//ast.Print(fileSet, astFile)
	//fmt.Println(newFileBuf.String())
	this.write(this.outfile, newFileBuf.String())
}
