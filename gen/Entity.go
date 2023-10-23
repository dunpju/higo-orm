package gen

import (
	"bytes"
	"fmt"
	"github.com/dunpju/higo-orm/gen/stubs"
	"github.com/dunpju/higo-utils/utils"
	"os"
	"strings"
)

const (
	entityStubFilename         = "entity.stub"
	entityPropertyStubFilename = "entityProperty.stub"
)

type property struct {
	isPrimaryKey      bool
	upperProperty     string
	propertyType      string
	tableField        string
	tableFieldComment string
}

func newProperty(isPrimaryKey bool, upperProperty, propertyType, tableField, tableFieldComment string) property {
	return property{isPrimaryKey: isPrimaryKey, upperProperty: upperProperty, propertyType: propertyType, tableField: tableField, tableFieldComment: tableFieldComment}
}

type Entity struct {
	stubContext         string
	table               Table
	primaryKey          string
	outDir              string
	entityFilename      string
	entityPackage       string
	outfile             string
	propertyString      []string
	imports             []string
	flags               []string
	properties          []property
	upperProperties     []string
	newFileBuf          *bytes.Buffer
	fieldMaxLen         int
	propertyTypeMaxLen  int
	upperPropertyMaxLen int
}

func (this *Entity) setProperties(properties []property) *Entity {
	this.properties = properties
	return this
}

func (this *Entity) setOutDir(outDir string) *Entity {
	this.outDir = outDir
	return this
}

func (this *Entity) setPackage(pkg string) *Entity {
	this.entityPackage = pkg
	return this
}

func (this *Entity) setFieldMaxLen(fieldMaxLen int) *Entity {
	this.fieldMaxLen = fieldMaxLen
	return this
}

func (this *Entity) setPropertyTypeMaxLen(propertyTypeMaxLen int) *Entity {
	this.propertyTypeMaxLen = propertyTypeMaxLen
	return this
}

func (this *Entity) setUpperPropertyMaxLen(upperPropertyMaxLen int) *Entity {
	this.upperPropertyMaxLen = upperPropertyMaxLen
	return this
}

func (this *Entity) setPrimaryKey(primaryKey string) *Entity {
	this.primaryKey = primaryKey
	return this
}

func (this *Entity) setTable(table Table) *Entity {
	this.table = table
	return this
}

func newEntity() *Entity {
	return &Entity{
		stubContext:     stubs.NewStub(entityStubFilename).Context(),
		entityFilename:  "Entity.go",
		imports:         make([]string, 0),
		flags:           make([]string, 0),
		propertyString:  make([]string, 0),
		upperProperties: make([]string, 0),
		newFileBuf:      bytes.NewBufferString(""),
	}
}

func (this *Entity) Write(p []byte) (n int, err error) {
	return this.newFileBuf.Write(p)
}

func (this *Entity) gen() {
	for _, p := range this.properties {
		if p.propertyType == "time.Time" {
			this.mergeImport(`"time"`)
		}
		blankFirst := LeftStrPad(" ", this.upperPropertyMaxLen-len(p.upperProperty), " ")
		blankSecond := LeftStrPad(" ", this.propertyTypeMaxLen-len(p.propertyType), " ")
		blankThree := LeftStrPad(" ", this.fieldMaxLen-len(p.tableField), " ")
		rowProperty := this.replaceRowProperty(p.upperProperty, blankFirst, p.propertyType, blankSecond, blankThree, p.tableField, p.tableFieldComment)
		this.mergeProperty(rowProperty)
	}
	this.replacePackage(this.entityPackage)
	this.replaceImport()
	this.replaceFields()
	this.replaceTableComment()
	this.replaceProperty()
	this.replacePrimaryKey(this.primaryKey)
	this.outfile = this.outDir + string(os.PathSeparator) + this.entityPackage + string(os.PathSeparator) + this.entityFilename
	if _, err := os.Stat(this.outfile); os.IsNotExist(err) {
		this.write(this.outfile, this.stubContext)
	} else {
		//this.oldAstEach(this.newAstEach())
	}
	fmt.Println(fmt.Sprintf("Entity IDE %s was created.", this.outfile))
}

func (this *Entity) write(file, fileContext string) {
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

func (this *Entity) mergeImport(ipt string) {
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

func (this *Entity) mergeProperty(rowProperty string) {
	has := false
	leftStrPad := LeftStrPad(rowProperty, 4, " ")
	for _, s := range this.propertyString {
		if s == leftStrPad {
			has = true
			break
		}
	}
	if !has {
		this.propertyString = append(this.propertyString, leftStrPad)
	}
}

func (this *Entity) replacePackage(pkg string) {
	this.stubContext = strings.Replace(this.stubContext, "%PACKAGE%", pkg, 1)
}

func (this *Entity) replaceImport() {
	imports := []string{
		LeftStrPad(`"github.com/dunpju/higo-orm/arm"`, 4, " "),
	}
	this.stubContext = strings.Replace(this.stubContext, "%IMPORT%", strings.Join(append(imports, this.imports...), "\n"), 1)
}

func (this *Entity) replaceFields() {
	flags := []string{
		LeftStrPad(`FlagDelete arm.Flag = iota + 1`, 4, " "),
		LeftStrPad(`FlagUpdate`, 4, " "),
	}
	this.stubContext = strings.Replace(this.stubContext, "%FLAGS%", strings.Join(append(flags, this.flags...), "\n"), 1)
}

func (this *Entity) replaceTableComment() {
	this.stubContext = strings.Replace(this.stubContext, "%TABLE_COMMENT%", this.table.Comment, 1)
}

func (this *Entity) replaceProperty() {
	editBlankSecond := LeftStrPad(" ", this.upperPropertyMaxLen-len("_edit"), " ")
	flagBlankSecond := LeftStrPad(" ", this.upperPropertyMaxLen-len("_flag"), " ")
	properties := []string{
		LeftStrPad(fmt.Sprintf("_edit%sbool", editBlankSecond), 4, " "),
		LeftStrPad(fmt.Sprintf("_flag%sarm.Flag", flagBlankSecond), 4, " "),
	}
	this.stubContext = strings.Replace(this.stubContext, "%PROPERTY%", strings.Join(append(properties, this.propertyString...), "\n"), 1)
}

func (this *Entity) replaceRowProperty(upperProperty, blankFirst, propertyType, blankSecond, blankThree, tableFields, tableFieldsComment string) string {
	stub := stubs.NewStub(entityPropertyStubFilename).Context()
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 1)
	stub = strings.Replace(stub, "%BLANK_FIRST%", blankFirst, 1)
	stub = strings.Replace(stub, "%PROPERTY_TYPE%", propertyType, 1)
	stub = strings.Replace(stub, "%BLANK_SECOND%", blankSecond, 1)
	stub = strings.Replace(stub, "%BLANK_THREE%", blankThree, 1)
	stub = strings.Replace(stub, "%TABLE_FIELDS%", tableFields, 1)
	stub = strings.Replace(stub, "%TABLE_FIELDS_COMMENT%", tableFieldsComment, 1)
	return stub
}

func (this *Entity) replacePrimaryKey(primaryKey string) {
	this.stubContext = strings.Replace(this.stubContext, "%PRIMARY_KEY%", primaryKey, 1)
}
