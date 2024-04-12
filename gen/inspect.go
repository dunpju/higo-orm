package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
)

type Inspect struct {
	alternativeAst *AlternativeAst
}

func NewInspect(alternativeAst *AlternativeAst) *Inspect {
	return &Inspect{alternativeAst: alternativeAst}
}

func (this *Inspect) Import(buffer *bytes.Buffer, n *ast.GenDecl) {
	if n.Tok.IsKeyword() && n.Tok.String() == token.IMPORT.String() {
		buffer.WriteString(fmt.Sprintf("%s ", token.IMPORT.String()))
		if n.Lparen.IsValid() {
			buffer.WriteString(fmt.Sprintf("%s\n", token.LPAREN.String()))
		}
		for _, spec := range n.Specs {
			buffer.WriteString(fmt.Sprintf("%s\n", LeftStrPad(spec.(*ast.ImportSpec).Path.Value, 4, " ")))
		}
		if n.Rparen.IsValid() {
			buffer.WriteString(fmt.Sprintf("%s\n", token.RPAREN.String()))
		}
		buffer.WriteString("\n")
	}
}

func (this *Inspect) Const(buffer *bytes.Buffer, n *ast.GenDecl) {
	if n.Tok.IsKeyword() && n.Tok.String() == token.CONST.String() {
		constNode := newConstNodeCollect()
		for _, spec := range n.Specs {
			constNode.append(spec.(*ast.ValueSpec))
		}
		for _, spec := range this.alternativeAst.constNode.Specs {
			constNode.append(spec.(*ast.ValueSpec))
		}
		upperPropertyMaxLen := 0
		for _, spec := range constNode.collect {
			for _, ident := range spec.Names {
				if upperPropertyMaxLen < len(ident.Name) {
					upperPropertyMaxLen = len(ident.Name)
				}
			}
		}
		buffer.WriteString(fmt.Sprintf("%s ", token.CONST.String()))
		if n.Lparen.IsValid() {
			buffer.WriteString(fmt.Sprintf("%s\n", token.LPAREN.String()))
		}
		for _, spec := range constNode.collect {
			blank := ""
			for _, ident := range spec.Names {
				buffer.WriteString(fmt.Sprintf("%s", LeftStrPad(ident.Name, 4, " ")))
				blank = fmt.Sprintf("%s", LeftStrPad(" ", upperPropertyMaxLen-len(ident.Name), " "))
				buffer.WriteString(blank)
			}
			selectorExpr := spec.Type.(*ast.SelectorExpr)
			buffer.WriteString(fmt.Sprintf("%s.%s ", selectorExpr.X.(*ast.Ident).Name, selectorExpr.Sel.Name))
			for _, expr := range spec.Values {
				buffer.WriteString(fmt.Sprintf("%s %s", token.ASSIGN, expr.(*ast.BasicLit).Value))
				buffer.WriteString(blank)
			}
			buffer.WriteString(fmt.Sprintf("%s%s%s", token.QUO, token.QUO, spec.Comment.Text()))
			pattern := `\n$`
			reg, _ := regexp.Compile(pattern)
			matched := reg.Match([]byte(buffer.String()))
			if !matched {
				buffer.WriteString("\n")
			}
		}
		if n.Rparen.IsValid() {
			buffer.WriteString(fmt.Sprintf("%s\n", token.RPAREN.String()))
		}
		buffer.WriteString("\n")
	}
}

func (this *Inspect) Type(buffer *bytes.Buffer, n *ast.GenDecl, hasStarExprArmModel bool) {
	if n.Specs != nil && len(n.Specs) > 0 {
		if typeSpec, ok := n.Specs[0].(*ast.TypeSpec); ok {
			structType, structTypeOk := typeSpec.Type.(*ast.StructType)
			if structTypeOk && typeSpec.Name.Obj.Kind.String() == token.TYPE.String() {
				if n.Doc.Text() != "" {
					buffer.WriteString(fmt.Sprintf("%s%s %s", token.QUO, token.QUO, n.Doc.Text()))
				}
				buffer.WriteString(fmt.Sprintf("%s ", n.Tok.String()))
				buffer.WriteString(fmt.Sprintf("%s ", typeSpec.Name.String()))
				buffer.WriteString(fmt.Sprintf("%s ", token.STRUCT.String()))
				buffer.WriteString(fmt.Sprintf("%s\n", token.LBRACE.String()))
				starExprList := newStarExprCollect()
				fieldsList := newFieldsListCollect()
				upperPropertyMaxLen := 0
				propertyTypeMaxLen := 0
				if typeSpec.Name.String() == modelStructName {
					for _, expr := range this.alternativeAst.starExprs {
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
							propertyTag := field.Tag.Value
							fieldsList.append(upperProperty, propertyType, propertyTag)
						}
					}
					for _, field := range this.alternativeAst.fieldsList {
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
						propertyTag := field.Tag.Value
						fieldsList.append(upperProperty, propertyType, propertyTag)
					}
				}
				for _, starExpr := range starExprList.collect {
					if starExprArmModel == starExpr {
						if hasStarExprArmModel {
							buffer.WriteString(LeftStrPad(fmt.Sprintf("%s%s\n", token.MUL, starExpr), 4, " "))
						}
					} else {
						buffer.WriteString(LeftStrPad(fmt.Sprintf("%s%s\n", token.MUL, starExpr), 4, " "))
					}
				}
				for _, fr := range fieldsList.collect {
					upperProperty := fmt.Sprintf("%s%s", fr.upperProperty, LeftStrPad(" ", upperPropertyMaxLen-len(fr.upperProperty), " "))
					propertyType := fmt.Sprintf("%s%s", fr.propertyType, LeftStrPad(" ", propertyTypeMaxLen-len(fr.propertyType), " "))
					buffer.WriteString(LeftStrPad(fmt.Sprintf("%s%s%s\n", upperProperty, propertyType, fr.propertyTag), 4, " "))
				}
				buffer.WriteString(fmt.Sprintf("%s\n", token.RBRACE.String()))
				buffer.WriteString("\n")
			}
		}
	}
}
