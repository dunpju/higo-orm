package gen

import (
	"bytes"
	"fmt"
	"github.com/dunpju/higo-orm/gen/stubs"
	"github.com/dunpju/higo-utils/utils"
	"os"
	"sort"
	"strings"
)

const (
	daoStubFilename          = "dao.stub"
	daoPropertyStubFilename  = "daoProperty.stub"
	daoInsertColumnsFilename = "daoInsertColumns.stub"
	daoValuesFilename        = "daoValues.stub"
	daoCaseFilename          = "daoCase.stub"
	daoWhenFilename          = "daoWhen.stub"
	daoCaseWhenFilename      = "daoCaseWhen.stub"
)

type ModelInfo struct {
	modelImport  string
	modelPackage string
}

func newModelInfo(modelImport string, modelPackage string) ModelInfo {
	return ModelInfo{modelImport: modelImport, modelPackage: modelPackage}
}

type EntityInfo struct {
	entityImport  string
	entityPackage string
}

func newEntityInfo(entityImport string, entityPackage string) EntityInfo {
	return EntityInfo{entityImport: entityImport, entityPackage: entityPackage}
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
	insertColumns       []string
	values              []string
	cases               []string
	whens               []string
	caseWhen            []string
	imports             []string
	flags               []string
	properties          []property
	newFileBuf          *bytes.Buffer
	fieldMaxLen         int
	propertyTypeMaxLen  int
	upperPropertyMaxLen int
	force               bool
}

func (this *Dao) setProperties(properties []property) *Dao {
	this.properties = properties
	return this
}

func (this *Dao) setOutDir(outDir string) *Dao {
	this.outDir = outDir
	return this
}

func (this *Dao) setDaoFilename(daoFilename string) *Dao {
	this.daoFilename = daoFilename
	return this
}

func (this *Dao) setModelInfo(modelInfo ModelInfo) *Dao {
	this.modelInfo = modelInfo
	return this
}

func (this *Dao) setEntityInfo(entityInfo EntityInfo) *Dao {
	this.entityInfo = entityInfo
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

func (this *Dao) setForce(force bool) *Dao {
	this.force = force
	return this
}

func newDao() *Dao {
	return &Dao{
		stubContext:    stubs.NewStub(daoStubFilename).Context(),
		imports:        make([]string, 0),
		flags:          make([]string, 0),
		propertyString: make([]string, 0),
		insertColumns:  make([]string, 0),
		values:         make([]string, 0),
		cases:          make([]string, 0),
		whens:          make([]string, 0),
		caseWhen:       make([]string, 0),
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
		if p.upperProperty == upperUpdateTime {
			rowUpdateTime = rowProperty
		}
		this.mergeProperty(rowProperty)
		this.mergeInsertColumns(p.upperProperty)
		this.mergeValues(p.upperProperty)
		this.mergeCase(p.lowerProperty, this.modelInfo.modelPackage, p.upperProperty)
		this.mergeWhen(p.lowerProperty, this.modelInfo.modelPackage, this.upperPrimaryKey, p.upperProperty)
		this.mergeCaseWhen(p.lowerProperty)
	}
	this.replacePackage(this.daoPackage)
	this.replaceImport()
	this.replaceModelProperty()
	this.replaceInsertColumns()
	this.replaceValues()
	this.replaceCases()
	this.replaceWhens()
	this.replaceCaseWhen()
	this.replacePrimaryKey(this.primaryKey)
	this.replaceUpperPrimaryKey(this.upperPrimaryKey)
	this.replaceModelPackage(this.modelInfo.modelPackage)
	this.replaceRowUpdateTime(rowUpdateTime)
	this.outfile = this.outDir + string(os.PathSeparator) + this.daoFilename
	if _, err := os.Stat(this.outfile); os.IsNotExist(err) {
		this.write(this.outfile, this.stubContext)
		fmt.Println(fmt.Sprintf("Dao IDE %s was created.", this.outfile))
	} else {
		if this.force {
			this.write(this.outfile, this.stubContext)
			fmt.Println(fmt.Sprintf("Dao IDE %s was forced updated.", this.outfile))
		} else {
			fmt.Println(fmt.Sprintf("Dao IDE %s was existent.", this.outfile))
		}
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

func (this *Dao) mergeProperty(rowProperty string) {
	has := false
	leftStrPad := LeftStrPad(rowProperty, 12, " ")
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

func (this *Dao) replaceRowProperty(upperProperty, blankFirst, tableFieldsComment string) string {
	stub := stubs.NewStub(daoPropertyStubFilename).Context()
	stub = strings.Replace(stub, "%MODEL_PACKAGE%", this.modelInfo.modelPackage, 1)
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 2)
	stub = strings.Replace(stub, "%BLANK_FIRST%", blankFirst, 1)
	stub = strings.Replace(stub, "%TABLE_FIELDS_COMMENT%", tableFieldsComment, 1)
	return stub
}

func (this *Dao) mergeInsertColumns(upperProperty string) {
	stub := stubs.NewStub(daoInsertColumnsFilename).Context()
	stub = strings.Replace(stub, "%MODEL_PACKAGE%", this.modelInfo.modelPackage, 1)
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 1)
	has := false
	leftStrPad := LeftStrPad(stub, 12, " ")
	if len(this.insertColumns) == 0 {
		leftStrPad = LeftStrPad(stub, 0, " ")
	}
	for _, s := range this.insertColumns {
		if s == leftStrPad {
			has = true
			break
		}
	}
	if !has {
		this.insertColumns = append(this.insertColumns, leftStrPad)
	}
}

func (this *Dao) replaceInsertColumns() {
	this.stubContext = strings.Replace(this.stubContext, "%INSERT_COLUMNS%", strings.Join(this.insertColumns, "\n"), 1)
}

func (this *Dao) mergeValues(upperProperty string) {
	stub := stubs.NewStub(daoValuesFilename).Context()
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 1)
	has := false
	leftStrPad := LeftStrPad(stub, 16, " ")
	if len(this.values) == 0 {
		leftStrPad = LeftStrPad(stub, 0, " ")
	}
	for _, s := range this.values {
		if s == leftStrPad {
			has = true
			break
		}
	}
	if !has {
		this.values = append(this.values, leftStrPad)
	}
}

func (this *Dao) replaceValues() {
	this.stubContext = strings.Replace(this.stubContext, "%VALUES%", strings.Join(this.values, "\n"), 1)
}

func (this *Dao) mergeCase(lowerProperty, modelPackage, upperProperty string) {
	stub := stubs.NewStub(daoCaseFilename).Context()
	stub = strings.Replace(stub, "%LOWER_PROPERTY%", lowerProperty, 1)
	stub = strings.Replace(stub, "%MODEL_PACKAGE%", modelPackage, 1)
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 1)
	has := false
	leftStrPad := LeftStrPad(stub, 8, " ")
	if len(this.cases) == 0 {
		leftStrPad = LeftStrPad(stub, 0, " ")
	}
	for _, s := range this.cases {
		if s == leftStrPad {
			has = true
			break
		}
	}
	if !has {
		this.cases = append(this.cases, leftStrPad)
	}
}

func (this *Dao) replaceCases() {
	this.stubContext = strings.Replace(this.stubContext, "%CASES%", strings.Join(this.cases, "\n"), 1)
}

func (this *Dao) mergeWhen(lowerProperty, modelPackage, upperPrimaryKey, upperProperty string) {
	stub := stubs.NewStub(daoWhenFilename).Context()
	stub = strings.Replace(stub, "%LOWER_PROPERTY%", lowerProperty, 1)
	stub = strings.Replace(stub, "%MODEL_PACKAGE%", modelPackage, 1)
	stub = strings.Replace(stub, "%UPPER_PRIMARY_KEY%", upperPrimaryKey, 2)
	stub = strings.Replace(stub, "%UPPER_PROPERTY%", upperProperty, 1)
	has := false
	leftStrPad := LeftStrPad(stub, 12, " ")
	if len(this.whens) == 0 {
		leftStrPad = LeftStrPad(stub, 0, " ")
	}
	for _, s := range this.whens {
		if s == leftStrPad {
			has = true
			break
		}
	}
	if !has {
		this.whens = append(this.whens, leftStrPad)
	}
}

func (this *Dao) replaceWhens() {
	this.stubContext = strings.Replace(this.stubContext, "%WHENS%", strings.Join(this.whens, "\n"), 1)
}

func (this *Dao) mergeCaseWhen(lowerProperty string) {
	stub := stubs.NewStub(daoCaseWhenFilename).Context()
	stub = strings.Replace(stub, "%LOWER_PROPERTY%", lowerProperty, 1)
	has := false
	leftStrPad := LeftStrPad(stub, 8, " ")
	if len(this.caseWhen) == 0 {
		leftStrPad = LeftStrPad(stub, 0, " ")
	}
	for _, s := range this.caseWhen {
		if s == leftStrPad {
			has = true
			break
		}
	}
	if !has {
		this.caseWhen = append(this.caseWhen, leftStrPad)
	}
}

func (this *Dao) replaceCaseWhen() {
	this.stubContext = strings.Replace(this.stubContext, "%CASE_WHEN%", strings.Join(this.caseWhen, "\n"), 1)
}

func (this *Dao) replacePackage(pkg string) {
	this.stubContext = strings.Replace(this.stubContext, "%PACKAGE%", pkg, 1)
}

func (this *Dao) replaceImport() {
	imports := []string{
		LeftStrPad(armImport, 4, " "),
		LeftStrPad(daoExceptionImport, 4, " "),
		LeftStrPad(himImport, 4, " "),
		LeftStrPad(gormImport, 4, " "),
		LeftStrPad(fmt.Sprintf(`"%s"`, this.entityInfo.entityImport), 4, " "),
		LeftStrPad(fmt.Sprintf(`"%s"`, this.modelInfo.modelImport), 4, " "),
	}
	sort.Strings(imports)
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
	if rowUpdateTime != "" {
		rowUpdateTime = "\n" + LeftStrPad(rowUpdateTime, 12, " ")
	}
	this.stubContext = strings.Replace(this.stubContext, "%ROW_UPDATE_TIME%", rowUpdateTime, 1)
}
