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
func (this *Entity) newAstEach() *AlternativeAst {
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
					if ok && typeSpec.Name.Obj.Kind.String() == token.TYPE.String() && typeSpec.Name.String() == entityStructName {
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
		}
		return true
	})
	//ast.Print(fileSet, astFile)
	return alternativeAst
}

func (this *Entity) oldAstEach(alternativeAst *AlternativeAst) {
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
	newFileBuf := bytes.NewBufferString("")
	ast.Inspect(astFile, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.File:
			newFileBuf.WriteString(fmt.Sprintf("package %s\n", n.Name.Name))
			newFileBuf.WriteString("\n")
		case *ast.GenDecl:
			if n.Tok.IsKeyword() && n.Tok.String() == token.IMPORT.String() {
				newFileBuf.WriteString(fmt.Sprintf("%s ", token.IMPORT.String()))
				if n.Lparen.IsValid() {
					newFileBuf.WriteString(fmt.Sprintf("%s\n", token.LPAREN.String()))
				}
				for _, spec := range n.Specs {
					newFileBuf.WriteString(fmt.Sprintf("%s\n", LeftStrPad(spec.(*ast.ImportSpec).Path.Value, 4, " ")))
				}
				if n.Rparen.IsValid() {
					newFileBuf.WriteString(fmt.Sprintf("%s\n", token.RPAREN.String()))
				}
				newFileBuf.WriteString("\n")
			} else if n.Tok.IsKeyword() && n.Tok.String() == token.CONST.String() {
				genDeclWrite := newGenDeclWrite()
				err = printer.Fprint(genDeclWrite, fileSet, n)
				if err != nil {
					panic(err)
				}
				genDeclWrite.buf.WriteString("\n")
				newFileBuf.WriteString(genDeclWrite.buf.String())
				newFileBuf.WriteString("\n")
			} else if n.Specs != nil && len(n.Specs) > 0 {
				if typeSpec, ok := n.Specs[0].(*ast.TypeSpec); ok {
					structType, structTypeOk := typeSpec.Type.(*ast.StructType)
					if structTypeOk && typeSpec.Name.Obj.Kind.String() == token.TYPE.String() {
						if n.Doc.Text() != "" {
							newFileBuf.WriteString(fmt.Sprintf("%s%s %s", token.QUO, token.QUO, n.Doc.Text()))
						}
						newFileBuf.WriteString(fmt.Sprintf("%s ", n.Tok.String()))
						newFileBuf.WriteString(fmt.Sprintf("%s ", typeSpec.Name.String()))
						newFileBuf.WriteString(fmt.Sprintf("%s ", token.STRUCT.String()))
						newFileBuf.WriteString(fmt.Sprintf("%s\n", token.LBRACE.String()))
						starExprList := newStarExprCollect()
						fieldsList := newFieldsListCollect()
						upperPropertyMaxLen := 0
						propertyTypeMaxLen := 0
						if typeSpec.Name.String() == entityStructName {
							for _, expr := range alternativeAst.starExprs {
								before := fmt.Sprintf("%s.%s", expr.X.(*ast.SelectorExpr).X.(*ast.Ident).String(),
									expr.X.(*ast.SelectorExpr).Sel.String())
								starExprList.append(before)
							}
							for _, field := range structType.Fields.List {
								starExpr, starExprOk := field.Type.(*ast.StarExpr)
								if starExprOk && len(field.Names) == 0 {
									after := fmt.Sprintf("%s.%s", starExpr.X.(*ast.SelectorExpr).X.(*ast.Ident).String(),
										starExpr.X.(*ast.SelectorExpr).Sel.String())
									starExprList.append(after)
								} else if len(field.Names) > 0 && !starExprOk {
									upperProperty := field.Names[0].Name
									if upperPropertyMaxLen < len(upperProperty) {
										upperPropertyMaxLen = len(upperProperty)
									}
									var propertyType string
									ident, identOK := field.Type.(*ast.Ident)
									if identOK {
										propertyType = ident.String()
									} else if selectorExpr, selectorExprOK := field.Type.(*ast.SelectorExpr); selectorExprOK {
										propertyType = fmt.Sprintf("%s.%s", selectorExpr.X.(*ast.Ident).String(), selectorExpr.Sel.String())
									} else if _, interfaceTypeOK := field.Type.(*ast.InterfaceType); interfaceTypeOK {
										propertyType = fmt.Sprintf("%s%s%s", token.INTERFACE, token.LBRACE, token.RBRACE)
									}
									if propertyTypeMaxLen < len(propertyType) {
										propertyTypeMaxLen = len(propertyType)
									}
									var propertyTag string
									if field.Tag != nil {
										propertyTag = field.Tag.Value
									}
									fieldsList.append(upperProperty, propertyType, propertyTag)
								}
							}
							for _, field := range alternativeAst.fieldsList {
								upperProperty := field.Names[0].Name
								if upperPropertyMaxLen < len(upperProperty) {
									upperPropertyMaxLen = len(upperProperty)
								}
								var propertyType string
								ident, identOK := field.Type.(*ast.Ident)
								if identOK {
									propertyType = ident.String()
								} else if selectorExpr, selectorExprOK := field.Type.(*ast.SelectorExpr); selectorExprOK {
									propertyType = fmt.Sprintf("%s.%s", selectorExpr.X.(*ast.Ident).String(), selectorExpr.Sel.String())
								}
								if propertyTypeMaxLen < len(propertyType) {
									propertyTypeMaxLen = len(propertyType)
								}
								var propertyTag string
								if field.Tag != nil {
									propertyTag = field.Tag.Value
								}
								fieldsList.append(upperProperty, propertyType, propertyTag)
							}
						}
						for _, starExpr := range starExprList.collect {
							newFileBuf.WriteString(LeftStrPad(fmt.Sprintf("%s%s\n", token.MUL, starExpr), 4, " "))
						}
						for _, fr := range fieldsList.collect {
							upperProperty := fmt.Sprintf("%s%s", fr.upperProperty, LeftStrPad(" ", upperPropertyMaxLen-len(fr.upperProperty), " "))
							propertyType := fmt.Sprintf("%s%s", fr.propertyType, LeftStrPad(" ", propertyTypeMaxLen-len(fr.propertyType), " "))
							newFileBuf.WriteString(LeftStrPad(fmt.Sprintf("%s%s%s\n", upperProperty, propertyType, fr.propertyTag), 4, " "))
						}
						newFileBuf.WriteString(fmt.Sprintf("%s\n", token.RBRACE.String()))
						newFileBuf.WriteString("\n")
					}
				}
			}
		case *ast.FuncDecl:
			isEntityFunc := false
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
																	if ident.Name == entityStructName && findProperty(selectorExpr.Sel.Name, this.upperProperties) {
																		funcList.append(newFnDecl(selectorExpr.Sel.Name, fileSet, n))
																		isEntityFunc = true
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
			if !isEntityFunc {
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
