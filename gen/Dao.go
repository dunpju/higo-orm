package gen

import (
	"bytes"
	"fmt"
	"github.com/dunpju/higo-orm/gen/stubs"
	"github.com/dunpju/higo-utils/utils"
	"os"
	"strings"
)

var UpperUpdateTime = "UpdateTime" // 大驼峰更新时间属性

const (
	daoStubFilename         = "dao.stub"
	daoPropertyStubFilename = "daoProperty.stub"
)

type ModelInfo struct {
	ipt          string
	modelPackage string
}

type EntityInfo struct {
	ipt           string
	entityPackage string
}

type Dao struct {
	stubContext         string
	table               Table
	primaryKey          string
	upperPrimaryKey     string
	outDir              string
	modelInfo           ModelInfo
	daoFilename         string
	daoPackage          string
	entityInfo          EntityInfo
	outfile             string
	propertyString      []string
	imports             []string
	flags               []string
	properties          []property
	newFileBuf          *bytes.Buffer
	fieldMaxLen         int
	propertyTypeMaxLen  int
	upperPropertyMaxLen int
}

func (this *Dao) setProperties(properties []property) *Dao {
	this.properties = properties
	return this
}

func (this *Dao) setOutDir(outDir string) *Dao {
	this.outDir = outDir
	return this
}

func (this *Dao) setModelInfo(modelInfo ModelInfo) *Dao {
	this.modelInfo = modelInfo
	return this
}

func (this *Dao) setPackage(pkg string) *Dao {
	this.daoPackage = pkg
	return this
}

func (this *Dao) setFieldMaxLen(fieldMaxLen int) *Dao {
	this.fieldMaxLen = fieldMaxLen
	return this
}

func (this *Dao) setPropertyTypeMaxLen(propertyTypeMaxLen int) *Dao {
	this.propertyTypeMaxLen = propertyTypeMaxLen
	return this
}

func (this *Dao) setUpperPropertyMaxLen(upperPropertyMaxLen int) *Dao {
	this.upperPropertyMaxLen = upperPropertyMaxLen
	return this
}

func (this *Dao) setPrimaryKey(primaryKey string) *Dao {
	this.primaryKey = primaryKey
	return this
}

func (this *Dao) setUpperPrimaryKey(upperPrimaryKey string) *Dao {
	this.upperPrimaryKey = upperPrimaryKey
	return this
}

func (this *Dao) setTable(table Table) *Dao {
	this.table = table
	return this
}

func newDao() *Dao {
	return &Dao{
		stubContext:    stubs.NewStub(daoStubFilename).Context(),
		imports:        make([]string, 0),
		flags:          make([]string, 0),
		propertyString: make([]string, 0),
		newFileBuf:     bytes.NewBufferString(""),
	}
}

func (this *Dao) Write(p []byte) (n int, err error) {
	return this.newFileBuf.Write(p)
}

func (this *Dao) gen() {
	var rowUpdateTime string
	for _, p := range this.properties {
		blankFirst := LeftStrPad(" ", (this.upperPropertyMaxLen-len(p.upperProperty))*2, " ")
		rowProperty := this.replaceRowProperty(p.upperProperty, blankFirst, p.tableFieldComment)
		if p.upperProperty == UpperUpdateTime {
			rowUpdateTime = rowProperty
		}
		this.mergeProperty(rowProperty)
	}
	this.replacePackage(this.daoPackage)
	this.replaceImport()
	this.replaceModelProperty()
	this.replacePrimaryKey(this.primaryKey)
	this.replaceUpperPrimaryKey(this.upperPrimaryKey)
	this.replaceModelPackage(this.modelInfo.modelPackage)
	this.replaceRowUpdateTime(rowUpdateTime)

	this.outfile = this.outDir + string(os.PathSeparator) + this.daoPackage + string(os.PathSeparator) + this.daoFilename
	if _, err := os.Stat(this.outfile); os.IsNotExist(err) {
		this.write(this.outfile, this.stubContext)
		fmt.Println(fmt.Sprintf("Dao IDE %s was created.", this.outfile))
	} else {
		fmt.Println(fmt.Sprintf("Dao IDE %s was existent.", this.outfile))
	}
}

func (this *Dao) write(file, fileContext string) {
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

func (this *Dao) replaceRowProperty(upperProperty, blankFirst, tableFieldsComment string) string {
	stub := stubs.NewStub(daoPropertyStubFilename).Context()
	stub = strings.Replace(stub, "%MODEL_PACKAGE%", this.modelInfo.modelPackage, 1)
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 2)
	stub = strings.Replace(stub, "%BLANK_FIRST%", blankFirst, 1)
	stub = strings.Replace(stub, "%TABLE_FIELDS_COMMENT%", tableFieldsComment, 1)
	return stub
}

func (this *Dao) mergeProperty(rowProperty string) {
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

func (this *Dao) replacePackage(pkg string) {
	this.stubContext = strings.Replace(this.stubContext, "%PACKAGE%", pkg, 1)
}

func (this *Dao) replaceImport() {
	imports := []string{
		LeftStrPad(`"github.com/dunpju/higo-orm/arm"`, 4, " "),
		LeftStrPad(`"github.com/dunpju/higo-orm/exception/DaoException"`, 4, " "),
		LeftStrPad(`"github.com/dunpju/higo-orm/him"`, 4, " "),
		LeftStrPad(`"gorm.io/gorm"`, 4, " "),
		LeftStrPad(this.entityInfo.ipt, 4, " "),
		LeftStrPad(this.modelInfo.ipt, 4, " "),
	}
	this.stubContext = strings.Replace(this.stubContext, "%IMPORT%", strings.Join(append(imports, this.imports...), "\n"), 1)
}

func (this *Dao) replaceModelProperty() {
	this.stubContext = strings.Replace(this.stubContext, "%MODEL_PROPERTY%", strings.Join(this.propertyString, "\n"), 1)
}

func (this *Dao) replacePrimaryKey(primaryKey string) {
	this.stubContext = strings.ReplaceAll(this.stubContext, "%PRIMARY_KEY%", primaryKey)
}

func (this *Dao) replaceUpperPrimaryKey(upperPrimaryKey string) {
	this.stubContext = strings.ReplaceAll(this.stubContext, "%UPPER_PRIMARY_KEY%", upperPrimaryKey)
}

func (this *Dao) replaceModelPackage(modelPackage string) {
	this.stubContext = strings.ReplaceAll(this.stubContext, "%MODEL_PACKAGE%", modelPackage)
}

func (this *Dao) replaceRowUpdateTime(rowUpdateTime string) {
	this.stubContext = strings.Replace(this.stubContext, "%ROW_UPDATE_TIME%", rowUpdateTime, 1)
}
