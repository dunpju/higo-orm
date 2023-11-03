package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dunpju/higo-orm/gen/stubs"
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-utils/utils"
	"github.com/dunpju/higo-utils/utils/dirutil"
	"github.com/dunpju/higo-utils/utils/stringutil"
	. "github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/spf13/cobra"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	table           string
	conn            string
	prefix          string
	out             string
	capitalBeganReg = regexp.MustCompile(`^[A-Z].*`) //匹配大写字母开头
)

const (
	allTable                      = "all"
	modelStubFilename             = "model.stub"
	modelFieldsStubFilename       = "modelFields.stub"
	modelPropertyStubFilename     = "modelProperty.stub"
	modelWithPropertyStubFilename = "modelWithProperty.stub"
	modelStructName               = "Model"
	starExprArmModel              = "arm.Model"
	yes                           = "yes"
	no                            = "no"
)

type YesNo string

func (this YesNo) Bool() bool {
	lower := strings.ToLower(string(this))
	if lower == "yes" {
		return true
	} else if lower == "no" {
		return false
	}
	panic(fmt.Errorf("undefined Constant"))
}

func initModel() {
	model.Flags().StringVarP(&table, "table", "t", "", "表名,all生成所有表模型")
	err := model.MarkFlagRequired("table")
	if err != nil {
		panic(err)
	}
	model.Flags().StringVarP(&conn, "conn", "c", "Default", "数据库连接,默认值:Default")
	model.Flags().StringVarP(&prefix, "prefix", "p", "", "数据库前缀,如:fm_")
	model.Flags().StringVarP(&out, "out", "o", "", "模型生成目录,如:app\\models")
	err = model.MarkFlagRequired("out")
	if err != nil {
		panic(err)
	}
	generator.AddCommand(model)
}

// go run .\bin\generator.go model --table=school --conn=Default --prefix=ts_ --out=app\models
// go run .\bin\generator.go model --table=all --conn=Default --prefix=ts_ --out=app\models
var model = &cobra.Command{
	Use:     "model",
	Short:   "模型构建工具",
	Long:    `模型构建工具`,
	Example: "model",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			isGenerateDao        YesNo = yes
			isGenerateEntity     YesNo
			confirmBeginGenerate YesNo
			isMatchCapitalBegan  string
			outEntityDir         string
			outDaoDir            string
		)
	loopDao:
		fmt.Print("Whether Generate Dao [yes|no] (default:yes):")
		n, err := fmt.Scanln(&isGenerateDao)
		if nil != err && n > 0 {
			panic(err)
		}
		if (yes != isGenerateDao && no != isGenerateDao) && n > 0 {
			goto loopDao
		}
		fmt.Printf("Choice Generate Dao: %s\n", isGenerateDao)
		if isGenerateDao.Bool() { // 确认构建dao
			if capitalBeganReg == nil {
				log.Fatalln("regexp err")
			}
			daoDir := "dao"
			isMatchCapitalBegan = capitalBeganReg.FindString(dirutil.Basename(out))
			if isMatchCapitalBegan != "" {
				daoDir = stringutil.Ucfirst(daoDir)
			}
			outDaoDir = dirutil.Dirname(out) + `\` + daoDir
			fmt.Printf("Confirm Output Directory Of Dao Default (%s)? Enter/Input: ", outDaoDir)
			n, err = fmt.Scanln(&outDaoDir)
			if nil != err && n > 0 {
				panic(err)
			}
			fmt.Printf("Confirmed Output Directory Of Dao: %s\n", outDaoDir)
			//确认构建dao，默认必须构建entity
			isGenerateEntity = yes
			goto loopChoiceGenerateEntity
		}
	loopEntity:
		fmt.Print("Whether Generate Entity [yes|no] (default:yes):")
		n, err = fmt.Scanln(&isGenerateEntity)
		if nil != err && n > 0 {
			panic(err)
		}
		if (yes != isGenerateEntity && no != isGenerateEntity) && n > 0 {
			goto loopEntity
		}
	loopChoiceGenerateEntity:
		fmt.Printf("Choice Generate Entity: %s\n", isGenerateEntity)
		if isGenerateEntity.Bool() { //确认构建entity
			entityDir := "entity"
			isMatchCapitalBegan = capitalBeganReg.FindString(dirutil.Basename(out))
			if isMatchCapitalBegan != "" {
				entityDir = stringutil.Ucfirst(entityDir)
			}
			outEntityDir = dirutil.Dirname(out) + `\` + entityDir
			fmt.Printf("Confirm Output Directory Of Entity Default (%s)? Enter/Input: ", outEntityDir)
			n, err = fmt.Scanln(&outEntityDir)
			if nil != err && n > 0 {
				panic(err)
			}
			fmt.Printf("Confirmed Output Directory Of Entity: %s\n", outEntityDir)
		}
		//确认开始构建
	loopConfirmBeginGenerate:
		fmt.Print("Whether Start Generate [yes|no] (default:yes):")
		n, err = fmt.Scanln(&confirmBeginGenerate)
		if (yes != confirmBeginGenerate && no != confirmBeginGenerate) && n > 0 {
			goto loopConfirmBeginGenerate
		}
		if (yes != confirmBeginGenerate) && n > 0 {
			goto loopDao
		}
		fmt.Print("Start Generate ......\n")
		db, err := him.DBConnect(conn)
		if err != nil {
			panic(err)
		}
		if prefix == "" {
			prefix = db.DBC().Prefix()
		}
		model := newModel(db, prefix, isGenerateDao, isGenerateEntity).setOutEntityDir(outEntityDir).setOutDaoDir(outDaoDir)
		if allTable == table {
			model.getTables()
		} else {
			model.getTable(table)
		}
		model.gen(out)
	},
}

type Model struct {
	db                  *him.DB
	tables              []Table
	prefix              string
	originalStubContext string
	stubContext         string
	filename            string
	modelFilename       string
	outfile             string
	imports             []string
	fields              []string
	tableComment        string
	properties          []string
	withProperty        []string
	upperProperties     []string
	newFileBuf          *bytes.Buffer
	isGenerateDao       YesNo
	isGenerateEntity    YesNo
	outEntityDir        string
	outDaoDir           string
}

func newModel(db *him.DB, prefix string, isGenerateDao, isGenerateEntity YesNo) *Model {
	return &Model{
		db:                  db,
		tables:              make([]Table, 0),
		prefix:              prefix,
		originalStubContext: stubs.NewStub(modelStubFilename).Context(),
		modelFilename:       "Model.go",
		isGenerateDao:       isGenerateDao,
		isGenerateEntity:    isGenerateEntity,
	}
}

func (this *Model) reset() {
	this.imports = make([]string, 0)
	this.fields = make([]string, 0)
	this.properties = make([]string, 0)
	this.withProperty = make([]string, 0)
	this.upperProperties = make([]string, 0)
	this.newFileBuf = bytes.NewBufferString("")
}

func (this *Model) setOutEntityDir(outEntityDir string) *Model {
	this.outEntityDir = outEntityDir
	return this
}

func (this *Model) setOutDaoDir(outDaoDir string) *Model {
	this.outDaoDir = outDaoDir
	return this
}

func (this *Model) gen(outDir string) {
	for _, t := range this.tables {
		this.reset()
		this.tableComment = t.Comment
		tableFields := this.getTableFields(t.Name)
		upperPropertyMaxLen := 0
		fieldMaxLen := 0
		propertyTypeMaxLen := 0
		primaryKey := ""
		upperPrimaryKey := ""
		isPrimaryKey := false
		properties := make([]property, 0)
		isBreak := false
	begin:
		for _, field := range tableFields {
			upperProperty := CamelCase(field.Field)
			if upperPropertyMaxLen < len(upperProperty) {
				upperPropertyMaxLen = len(upperProperty)
			}
			if fieldMaxLen < len(field.Field) {
				fieldMaxLen = len(field.Field)
			}
			propertyType := convertFiledType(field)
			if propertyTypeMaxLen < len(propertyType) {
				propertyTypeMaxLen = len(propertyType)
			}
			if !isBreak {
				continue
			}
			if field.Key == "PRI" {
				upperPrimaryKey = upperProperty
				primaryKey = utils.String.Lcfirst(upperProperty)
				isPrimaryKey = true
			}
			if propertyType == "time.Time" {
				this.mergeImport(`"time"`)
			}
			this.appendProperty(upperProperty)
			blankFirst := LeftStrPad(" ", upperPropertyMaxLen-len(upperProperty), " ")
			rowField := this.replaceRowField(upperProperty, blankFirst, field.Field, blankFirst, field.Comment)
			this.mergeFields(rowField)
			blankSecond := LeftStrPad(" ", propertyTypeMaxLen-len(propertyType), " ")
			blankThree := LeftStrPad(" ", fieldMaxLen-len(field.Field), " ")
			blankFour := LeftStrPad(" ", fieldMaxLen-len(field.Field), " ")
			rowProperty := this.replaceRowProperty(upperProperty, blankFirst, propertyType, blankSecond, blankThree, blankFour, field.Field, field.Comment)
			this.mergeProperty(rowProperty)
			rowWithProperty := this.replaceRowWithProperty(upperProperty, utils.String.Lcfirst(upperProperty), propertyType, field.Comment)
			this.mergeWithProperty(rowWithProperty)
			properties = append(properties, newProperty(isPrimaryKey, upperProperty, propertyType, field.Field, field.Comment))
		}
		if fieldMaxLen > 0 {
			if !isBreak {
				isBreak = true
				goto begin
			}
		}
		this.stubContext = this.originalStubContext
		modelPackage := CamelCase(strings.Replace(t.Name, this.prefix, "", 1))
		this.outfile = outDir + string(os.PathSeparator) + modelPackage + string(os.PathSeparator) + this.modelFilename
		this.replacePackage(modelPackage)
		this.replaceImport()
		this.replaceFields()
		this.replaceTableComment()
		this.replaceProperty()
		this.replaceTableName(t.Name)
		this.replaceUpperPrimaryKey(upperPrimaryKey)
		this.replaceWithProperty()
		if _, err := os.Stat(this.outfile); os.IsNotExist(err) {
			this.write(this.outfile, this.stubContext)
		} else {
			this.oldAstEach(this.newAstEach())
		}
		fmt.Println(fmt.Sprintf("Model IDE %s was created.", this.outfile))
		entityPackage := fmt.Sprintf("%sEntity", modelPackage)
		if this.isGenerateDao.Bool() {
			newEntity().
				setOutDir(this.outEntityDir).
				setPackage(entityPackage).
				setTable(t).
				setUpperPrimaryKey(upperPrimaryKey).
				setProperties(properties).
				setFieldMaxLen(fieldMaxLen).
				setPropertyTypeMaxLen(propertyTypeMaxLen).
				setUpperPropertyMaxLen(upperPropertyMaxLen).
				gen()
			goMod := GetModInfo()
			pwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			childPath := GetGoModChildPath(pwd)
			var childPathStr string
			if len(childPath) > 0 {
				childPathStr = fmt.Sprintf("/%s", strings.Join(childPath, "/"))
			}
			modelImport := goMod.Module.Path + fmt.Sprintf("/%s%s", childPathStr, strings.ReplaceAll(utils.Dir.Dirname(this.outfile), "\\", "/"))
			entityImport := goMod.Module.Path + fmt.Sprintf("/%s%s", childPathStr, strings.ReplaceAll(fmt.Sprintf("%s/%s", this.outEntityDir, entityPackage), "\\", "/"))
			daoFilename := fmt.Sprintf("%sDao.go", modelPackage)
			newDao().
				setOutDir(this.outDaoDir).
				setDaoFilename(daoFilename).
				setPackage(utils.Dir.Basename(this.outDaoDir)).
				setTable(t).
				setModelInfo(newModelInfo(modelImport, modelPackage)).
				setEntityInfo(newEntityInfo(entityImport, entityPackage)).
				setUpperPrimaryKey(upperPrimaryKey).
				setPrimaryKey(primaryKey).
				setProperties(properties).
				setFieldMaxLen(fieldMaxLen).
				setPropertyTypeMaxLen(propertyTypeMaxLen).
				setUpperPropertyMaxLen(upperPropertyMaxLen).
				gen()
		} else if this.isGenerateEntity.Bool() {
			newEntity().
				setOutDir(this.outEntityDir).
				setPackage(entityPackage).
				setTable(t).
				setUpperPrimaryKey(upperPrimaryKey).
				setProperties(properties).
				setFieldMaxLen(fieldMaxLen).
				setPropertyTypeMaxLen(propertyTypeMaxLen).
				setUpperPropertyMaxLen(upperPropertyMaxLen).
				gen()
		}
	}
}

func (this *Model) write(file, fileContext string) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		utils.Dir.Mkdir(file, os.ModePerm)
	}
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write([]byte(fileContext))
	if err != nil {
		panic(err)
	}
}

func (this *Model) appendProperty(upperProperty string) {
	if !this.findProperty(upperProperty) {
		this.upperProperties = append(this.upperProperties, upperProperty)
	}
}

func (this *Model) findProperty(upperProperty string) bool {
	for _, s := range this.upperProperties {
		if s == upperProperty {
			return true
		}
	}
	return false
}

func (this *Model) mergeImport(ipt string) {
	has := false
	for _, s := range this.imports {
		if s == LeftStrPad(ipt, 4, " ") {
			has = true
			break
		}
	}
	if !has {
		this.imports = append(this.imports, LeftStrPad(ipt, 4, " "))
	}
}

func (this *Model) mergeFields(rawField string) {
	has := false
	leftStrPad := LeftStrPad(rawField, 4, " ")
	for _, s := range this.fields {
		if s == leftStrPad {
			has = true
			break
		}
	}
	if !has {
		this.fields = append(this.fields, leftStrPad)
	}
}

func (this *Model) mergeProperty(rowProperty string) {
	has := false
	leftStrPad := LeftStrPad(rowProperty, 4, " ")
	for _, s := range this.properties {
		if s == leftStrPad {
			has = true
			break
		}
	}
	if !has {
		this.properties = append(this.properties, leftStrPad)
	}
}

func (this *Model) mergeWithProperty(rowWithProperty string) {
	has := false
	for _, s := range this.withProperty {
		if s == rowWithProperty {
			has = true
			break
		}
	}
	if !has {
		this.withProperty = append(this.withProperty, rowWithProperty)
	}
}

func (this *Model) replaceRowField(upperProperty, blankFirst, tableFields, blankSecond, tableFieldsComment string) string {
	stub := stubs.NewStub(modelFieldsStubFilename).Context()
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 1)
	stub = strings.Replace(stub, "%BLANK_FIRST%", blankFirst, 1)
	stub = strings.Replace(stub, "%TABLE_FIELDS%", tableFields, 1)
	stub = strings.Replace(stub, "%BLANK_SECOND%", blankSecond, 1)
	stub = strings.Replace(stub, "%TABLE_FIELDS_COMMENT%", tableFieldsComment, 1)
	return stub
}

func (this *Model) replaceRowProperty(upperProperty, blankFirst, propertyType, blankSecond, blankThree, blankFour, tableFields, tableFieldsComment string) string {
	stub := stubs.NewStub(modelPropertyStubFilename).Context()
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 1)
	stub = strings.Replace(stub, "%BLANK_FIRST%", blankFirst, 1)
	stub = strings.Replace(stub, "%PROPERTY_TYPE%", propertyType, 1)
	stub = strings.Replace(stub, "%BLANK_SECOND%", blankSecond, 1)
	stub = strings.Replace(stub, "%BLANK_THREE%", blankThree, 1)
	stub = strings.Replace(stub, "%BLANK_FOUR%", blankFour, 1)
	stub = strings.Replace(stub, "%TABLE_FIELDS%", tableFields, 2)
	stub = strings.Replace(stub, "%TABLE_FIELDS_COMMENT%", tableFieldsComment, 1)
	return stub
}

func (this *Model) replaceRowWithProperty(upperProperty, lowerProperty, propertyType, tableFieldsComment string) string {
	stub := stubs.NewStub(modelWithPropertyStubFilename).Context()
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 3)
	stub = strings.Replace(stub, "%LOWER_PROPERTY%", lowerProperty, 2)
	stub = strings.Replace(stub, "%PROPERTY_TYPE%", propertyType, 1)
	stub = strings.Replace(stub, "%TABLE_FIELDS_COMMENT%", tableFieldsComment, 1)
	return stub
}

func (this *Model) replacePackage(pkg string) {
	this.stubContext = strings.Replace(this.stubContext, "%PACKAGE%", pkg, 1)
}

func (this *Model) replaceImport() {
	imports := []string{
		LeftStrPad(`"github.com/dunpju/higo-orm/arm"`, 4, " "),
		LeftStrPad(`"github.com/dunpju/higo-orm/him"`, 4, " "),
	}
	this.stubContext = strings.Replace(this.stubContext, "%IMPORT%", strings.Join(append(imports, this.imports...), "\n"), 1)
}

func (this *Model) replaceFields() {
	this.stubContext = strings.Replace(this.stubContext, "%FIELDS%", strings.Join(this.fields, "\n"), 1)
}

func (this *Model) replaceTableComment() {
	this.stubContext = strings.Replace(this.stubContext, "%TABLE_COMMENT%", this.tableComment, 1)
}

func (this *Model) replaceProperty() {
	this.stubContext = strings.Replace(this.stubContext, "%PROPERTY%", strings.Join(this.properties, "\n"), 1)
}

func (this *Model) replaceTableName(tableName string) {
	this.stubContext = strings.Replace(this.stubContext, "%TABLE_NAME%", tableName, 1)
}

func (this *Model) replaceUpperPrimaryKey(upperPrimaryKey string) {
	this.stubContext = strings.Replace(this.stubContext, "%UPPER_PRIMARY_KEY%", upperPrimaryKey, 1)
}

func (this *Model) replaceWithProperty() {
	this.stubContext = strings.Replace(this.stubContext, "%WITH_PROPERTY%", strings.Join(this.withProperty, "\n\n"), 1)
}

// getTables 获取数据库所有表
func (this *Model) getTables() {
	gormDB := this.db.Raw(fmt.Sprintf(`SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='%s' AND TABLE_NAME LIKE '%s%%'`, this.db.DBC().Database(), this.prefix)).Get(&this.tables)
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
}

// getTable 获取数据库表
func (this *Model) getTable(table string) {
	gormDB := this.db.Raw(fmt.Sprintf(`SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='%s' AND TABLE_NAME = '%s'`, this.db.DBC().Database(), table)).Get(&this.tables)
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
}

// getTableFields 获取表所有字段信息
func (this *Model) getTableFields(tableName string) []TableField {
	var fields []TableField
	gormDB := this.db.Raw(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", tableName)).Get(&fields)
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
	return fields
}

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
																	if ident.Name == modelStructName && this.findProperty(selectorExpr.Sel.Name) {
																		alternativeAst.funcList = append(alternativeAst.funcList, newFnDecl(selectorExpr.Sel.Name, n))
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
									hasStarExprArmModel = true
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
				constNode := newConstNodeCollect()
				for _, spec := range alternativeAst.constNode.Specs {
					constNode.append(spec.(*ast.ValueSpec))
				}
				for _, spec := range n.Specs {
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
				newFileBuf.WriteString(fmt.Sprintf("%s ", token.CONST.String()))
				if n.Lparen.IsValid() {
					newFileBuf.WriteString(fmt.Sprintf("%s\n", token.LPAREN.String()))
				}
				for _, spec := range constNode.collect {
					blank := ""
					for _, ident := range spec.Names {
						newFileBuf.WriteString(fmt.Sprintf("%s", LeftStrPad(ident.Name, 4, " ")))
						blank = fmt.Sprintf("%s", LeftStrPad(" ", upperPropertyMaxLen-len(ident.Name), " "))
						newFileBuf.WriteString(blank)
					}
					selectorExpr := spec.Type.(*ast.SelectorExpr)
					newFileBuf.WriteString(fmt.Sprintf("%s.%s ", selectorExpr.X.(*ast.Ident).Name, selectorExpr.Sel.Name))
					for _, expr := range spec.Values {
						newFileBuf.WriteString(fmt.Sprintf("%s %s", token.ASSIGN, expr.(*ast.BasicLit).Value))
						newFileBuf.WriteString(blank)
					}
					newFileBuf.WriteString(fmt.Sprintf("%s%s%s", token.QUO, token.QUO, spec.Comment.Text()))
					pattern := `\n$`
					reg, _ := regexp.Compile(pattern)
					matched := reg.Match([]byte(newFileBuf.String()))
					if !matched {
						newFileBuf.WriteString("\n")
					}
				}
				if n.Rparen.IsValid() {
					newFileBuf.WriteString(fmt.Sprintf("%s\n", token.RPAREN.String()))
				}
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
						if typeSpec.Name.String() == modelStructName {
							for _, expr := range alternativeAst.starExprs {
								before := fmt.Sprintf("%s.%s", expr.X.(*ast.SelectorExpr).X.(*ast.Ident).String(),
									expr.X.(*ast.SelectorExpr).Sel.String())
								starExprList.append(before)
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
								propertyTag := field.Tag.Value
								fieldsList.append(upperProperty, propertyType, propertyTag)
							}
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
						for _, starExpr := range starExprList.collect {
							if starExprArmModel == starExpr {
								if hasStarExprArmModel {
									newFileBuf.WriteString(LeftStrPad(fmt.Sprintf("%s%s\n", token.MUL, starExpr), 4, " "))
								}
							} else {
								newFileBuf.WriteString(LeftStrPad(fmt.Sprintf("%s%s\n", token.MUL, starExpr), 4, " "))
							}
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
																	if ident.Name == modelStructName && this.findProperty(selectorExpr.Sel.Name) {
																		funcList.append(newFnDecl(selectorExpr.Sel.Name, n))
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
		err = printer.Fprint(funcDeclWrite, fileSet, fd.Fd)
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

func (this *Model) Write(p []byte) (n int, err error) {
	return this.newFileBuf.Write(p)
}

type Table struct {
	Name    string `gorm:"column:Name" json:"name"`
	Comment string `gorm:"column:Comment" json:"comment"`
}

type StructField struct {
	FieldName         string
	FieldType         string
	TableFieldName    string
	TableFieldComment string
}

type TableField struct {
	Field      string `gorm:"column:Field"`
	Type       string `gorm:"column:Type"`
	Null       string `gorm:"column:Null"` //非空 YES/NO
	Key        string `gorm:"column:Key"`
	Default    string `gorm:"column:Default"`
	Extra      string `gorm:"column:Extra"`
	Privileges string `gorm:"column:Privileges"`
	Comment    string `gorm:"column:Comment"`
}

// 转换字段类型
func convertFiledType(field TableField) string {
	types := strings.Split(field.Type, "(")
	switch types[0] {
	case "int":
		return "int"
	case "integer":
		return "int"
	case "mediumint":
		return "int"
	case "bit":
		return "int"
	case "year":
		return "int"
	case "smallint":
		return "int"
	case "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "decimal":
		return "float32"
	case "double":
		return "float32"
	case "float":
		return "float32"
	case "real":
		return "float32"
	case "numeric":
		return "float32"
	case "timestamp":
		return "time.Time"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	case "binary":
		return "[]byte"
	case "varchar":
		return "string"
	default:
		return "interface{}"
	}
}

// LeftStrPad
// input string 原字符串
// padLength int 规定补齐后的字符串位数
// padString string 自定义填充字符串
func LeftStrPad(input string, padLength int, padString string) string {
	output := ""
	for i := 1; i <= padLength; i++ {
		output += padString
	}
	return output + input
}

func GetModInfo() *GoMod {
	cmd := exec.Command("go", "mod", "edit", "-json")
	buffer := bytes.NewBufferString("")
	cmd.Stdout = buffer
	cmd.Stderr = buffer

	if err := cmd.Run(); err != nil {
		panic(err)
	}
	goMod := &GoMod{}
	err := json.Unmarshal(buffer.Bytes(), &goMod)
	if err != nil {
		panic(err)
	}
	return goMod
}

type GoMod struct {
	Module  Module
	Go      string
	Require []Require
	Exclude []Module
}

type Module struct {
	Path    string
	Version string
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

func Dirslice(path string) []string {
	pathSeparator := string(os.PathSeparator)
	paths := strings.Split(path, pathSeparator)
	if len(paths) == 1 {
		re := regexp.MustCompile("/")
		if re.Match([]byte(paths[0])) {
			paths = strings.Split(path, pathSeparator)
		}
	}
	return paths
}
