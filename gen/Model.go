package gen

import (
	"fmt"
	"github.com/dunpju/higo-orm/gen/stubs"
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-utils/utils"
	. "github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/spf13/cobra"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

var (
	table  string
	conn   string
	prefix string
	out    string
)

const (
	allTable                      = "all"
	modelStubFilename             = "model.stub"
	modelFieldsStubFilename       = "modelFields.stub"
	modelPropertyStubFilename     = "modelProperty.stub"
	modelWithPropertyStubFilename = "modelWithProperty.stub"
)

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
var model = &cobra.Command{
	Use:     "model",
	Short:   "模型构建工具",
	Long:    `模型构建工具`,
	Example: "model",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := him.DBConnect(conn)
		if err != nil {
			panic(err)
		}
		if prefix == "" {
			prefix = db.DBC().Prefix()
		}
		model := newModel(db, prefix)
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
	outfile             string
	imports             []string
	fields              []string
	tableComment        string
	properties          []string
	withProperty        []string
}

func newModel(db *him.DB, prefix string) *Model {
	return &Model{
		db:                  db,
		tables:              make([]Table, 0),
		prefix:              prefix,
		originalStubContext: stubs.NewStub(modelStubFilename).Context(),
		filename:            "Model.go",
		imports:             make([]string, 0),
		fields:              make([]string, 0),
		properties:          make([]string, 0),
		withProperty:        make([]string, 0),
	}
}

func (this *Model) gen(out string) {
	for _, t := range this.tables {
		this.tableComment = t.Comment
		tableFields := this.getTableFields(t.Name)
		upperPropertyMaxLen := 0
		fieldMaxLen := 0
		propertyTypeMaxLen := 0
		primaryKey := ""
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
				primaryKey = upperProperty
			}
			if propertyType == "time.Time" {
				this.mergeImport(`"time"`)
			}
			blankFirst := LeftStrPad(" ", upperPropertyMaxLen-len(upperProperty), " ")
			rawField := this.replaceRawField(upperProperty, blankFirst, field.Field, blankFirst, field.Comment)
			this.mergeFields(rawField)
			blankSecond := LeftStrPad(" ", propertyTypeMaxLen-len(propertyType), " ")
			blankThree := LeftStrPad(" ", fieldMaxLen-len(field.Field), " ")
			blankFour := LeftStrPad(" ", fieldMaxLen-len(field.Field), " ")
			rawProperty := this.replaceRawProperty(upperProperty, blankFirst, propertyType, blankSecond, blankThree, blankFour, field.Field, field.Comment)
			this.mergeProperty(rawProperty)
			rawWithProperty := this.replaceRawWithProperty(upperProperty, utils.String.Lcfirst(upperProperty), propertyType, field.Comment)
			this.mergeWithProperty(rawWithProperty)
		}
		if fieldMaxLen > 0 {
			if !isBreak {
				isBreak = true
				goto begin
			}
		}
		this.stubContext = this.originalStubContext
		pkg := CamelCase(strings.Replace(t.Name, this.prefix, "", 1))
		this.replacePackage(pkg)
		this.replaceImport()
		this.replaceFields()
		this.replaceTableComment()
		this.replaceProperty()
		this.replaceTableName(t.Name)
		this.replacePrimaryKey(primaryKey)
		this.replaceWithProperty()
		this.astFindStructNode()
		this.write(out + string(os.PathSeparator) + pkg + string(os.PathSeparator) + this.filename)
	}
}

func (this *Model) write(file string) {
	this.outfile = file
	if _, err := os.Stat(file); os.IsNotExist(err) {
		utils.Dir.Mkdir(file, os.ModePerm)
	}
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write([]byte(this.stubContext))
	if err != nil {
		panic(err)
	}
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

func (this *Model) mergeProperty(rawProperty string) {
	has := false
	leftStrPad := LeftStrPad(rawProperty, 4, " ")
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

func (this *Model) mergeWithProperty(rawWithProperty string) {
	has := false
	for _, s := range this.withProperty {
		if s == rawWithProperty {
			has = true
			break
		}
	}
	if !has {
		this.withProperty = append(this.withProperty, rawWithProperty)
	}
}

func (this *Model) replaceRawField(upperProperty, blankFirst, tableFields, blankSecond, tableFieldsComment string) string {
	stub := stubs.NewStub(modelFieldsStubFilename).Context()
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 1)
	stub = strings.Replace(stub, "%BLANK_FIRST%", blankFirst, 1)
	stub = strings.Replace(stub, "%TABLE_FIELDS%", tableFields, 1)
	stub = strings.Replace(stub, "%BLANK_SECOND%", blankSecond, 1)
	stub = strings.Replace(stub, "%TABLE_FIELDS_COMMENT%", tableFieldsComment, 1)
	return stub
}

func (this *Model) replaceRawProperty(upperProperty, blankFirst, propertyType, blankSecond, blankThree, blankFour, tableFields, tableFieldsComment string) string {
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

func (this *Model) replaceRawWithProperty(upperProperty, lowerProperty, propertyType, tableFieldsComment string) string {
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

func (this *Model) replacePrimaryKey(primaryKey string) {
	this.stubContext = strings.Replace(this.stubContext, "%PRIMARY_KEY%", primaryKey, 1)
}

func (this *Model) replaceWithProperty() {
	this.stubContext = strings.Replace(this.stubContext, "%WITH_PROPERTY%", strings.Join(this.withProperty, "\n\n"), 1)
}

// GetTables 获取数据库所有表
func (this *Model) getTables() {
	gormDB := this.db.Raw(fmt.Sprintf(`SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='%s' AND TABLE_NAME LIKE '%s%%'`, this.db.DBC().Database(), this.prefix)).Get(&this.tables)
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
}

// GetTable 获取数据库表
func (this *Model) getTable(table string) {
	gormDB := this.db.Raw(fmt.Sprintf(`SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='%s' AND TABLE_NAME = '%s'`, this.db.DBC().Database(), table)).Get(&this.tables)
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
}

// GetTableFields 获取表所有字段信息
func (this *Model) getTableFields(tableName string) []TableField {
	var fields []TableField
	gormDB := this.db.Raw(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", tableName)).Get(&fields)
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
	return fields
}

// astFindStructNode 找到Struct Node
func (this *Model) astFindStructNode() {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, this.outfile, this.stubContext, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	var structNode *ast.GenDecl
	ast.Inspect(file, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.GenDecl:
			if n.Specs != nil && len(n.Specs) > 0 {
				if typeSpec, ok := n.Specs[0].(*ast.TypeSpec); ok {
					if typeSpec.Name.String() == "Model" {
						fmt.Println(typeSpec.Doc)
						structNode = n //找到struct node
					}
				}
			}
		}
		return true
	})
	fmt.Println(structNode)
	//ast.Print(fileSet, file)
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
