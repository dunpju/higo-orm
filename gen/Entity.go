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
	entityStructName           = "Entity"
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
	upperPrimaryKey     string
	outDir              string
	entityFilename      string
	entityPackage       string
	outfile             string
	propertyString      []string
	imports             []string
	flags               []string
	upperProperties     []string
	properties          []property
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

func (this *Entity) setUpperPrimaryKey(upperPrimaryKey string) *Entity {
	this.upperPrimaryKey = upperPrimaryKey
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
		upperProperties: make([]string, 0),
		propertyString:  make([]string, 0),
		newFileBuf:      bytes.NewBufferString(""),
	}
}

func (this *Entity) Write(p []byte) (n int, err error) {
	return this.newFileBuf.Write(p)
}

func (this *Entity) gen() {
	var (
		rowTimeNow string
		createTime string
		updateTime string
	)
	for _, p := range this.properties {
		if p.propertyType == timeImport {
			this.mergeImport(`"time"`)
		}
		blankFirst := LeftStrPad(" ", this.upperPropertyMaxLen-len(p.upperProperty), " ")
		blankSecond := LeftStrPad(" ", this.propertyTypeMaxLen-len(p.propertyType), " ")
		blankThree := LeftStrPad(" ", this.fieldMaxLen-len(p.tableField), " ")
		rowProperty := this.replaceRowProperty(p.upperProperty, blankFirst, p.propertyType, blankSecond, blankThree, p.tableField, p.tableFieldComment)
		if p.upperProperty == UpperCreateTime {
			rowTimeNow = timeNow
			createTime = p.upperProperty
		} else if p.upperProperty == UpperUpdateTime {
			rowTimeNow = timeNow
			updateTime = p.upperProperty
		}
		this.mergeProperty(rowProperty)
	}
	this.replacePackage(this.entityPackage)
	this.replaceImport()
	this.replaceFlags()
	this.replaceTableComment()
	this.replaceProperty()
	this.replaceTimeNow(rowTimeNow)
	this.replaceCreateUpdateTime(createTime, updateTime)
	this.replaceUpperPrimaryKey(this.upperPrimaryKey)
	this.outfile = this.outDir + string(os.PathSeparator) + this.entityPackage + string(os.PathSeparator) + this.entityFilename
	if _, err := os.Stat(this.outfile); os.IsNotExist(err) {
		this.write(this.outfile, this.stubContext)
		fmt.Println(fmt.Sprintf("Entity IDE %s was created.", this.outfile))
	} else {
		this.oldAstEach(this.newAstEach())
		fmt.Println(fmt.Sprintf("Entity IDE %s was updated.", this.outfile))
	}
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
		LeftStrPad(armImport, 4, " "),
	}
	this.stubContext = strings.Replace(this.stubContext, "%IMPORT%", strings.Join(append(imports, this.imports...), "\n"), 1)
}

func (this *Entity) replaceFlags() {
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

func (this *Entity) replaceUpperPrimaryKey(upperPrimaryKey string) {
	this.stubContext = strings.Replace(this.stubContext, "%UPPER_PRIMARY_KEY%", upperPrimaryKey, 1)
}

func (this *Entity) replaceTimeNow(timeNow string) {
	if timeNow != "" {
		timeNow = "\n" + LeftStrPad(fmt.Sprintf("tn := %s", timeNow), 4, " ")
	}
	this.stubContext = strings.Replace(this.stubContext, "%TIME_NOW%", timeNow, 1)
}

func (this *Entity) replaceCreateUpdateTime(createTime, updateTime string) {
	var row string
	if createTime != "" {
		row = fmt.Sprintf("%s: tn", createTime)
	}
	if updateTime != "" {
		if row != "" {
			row += ", "
		}
		row += fmt.Sprintf("%s: tn", updateTime)
	}
	this.stubContext = strings.Replace(this.stubContext, "%CREATE_UPDATE_TIME%", row, 1)
}
