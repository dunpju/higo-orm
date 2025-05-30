package gen

import (
	"bytes"
	"fmt"
	"github.com/dunpju/higo-orm/gen/stubs"
	"github.com/dunpju/higo-orm/him"
	"github.com/dunpju/higo-utils/utils"
	"github.com/dunpju/higo-utils/utils/dirutil"
	"github.com/dunpju/higo-utils/utils/stringutil"
	"github.com/spf13/cobra"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	table           string
	conn            string
	prefix          string
	out             string
	force           bool
	forceEntity     bool
	forceDao        bool
	all             bool
	upperCreateTime string
	upperUpdateTime string
	capitalBeganReg = regexp.MustCompile(`^[A-Z].*`) //匹配大写字母开头
)

const (
	allTable                      = "all"
	modelStubFilename             = "model.stub"
	modelFieldsStubFilename       = "modelFields.stub"
	modelPropertyStubFilename     = "modelProperty.stub"
	modelWithPropertyStubFilename = "modelWithProperty.stub"
	modelEqualStubFilename        = "modelEqual.stub"
	modelStructName               = "Model"
	starExprArmModel              = "arm.Model"
	Yes                           = "yes"
	No                            = "no"
)

func InitModel() {
	ModelCommand.Flags().StringVarP(&table, "table", "t", "", "表名")
	ModelCommand.Flags().StringVarP(&conn, "conn", "c", "Default", "数据库连接")
	ModelCommand.Flags().StringVarP(&prefix, "prefix", "p", "", "数据表前缀,如:fm_")
	ModelCommand.Flags().StringVarP(&out, "out", "o", "", "模型生成目录,如:app\\models")
	err := ModelCommand.MarkFlagRequired("out")
	if err != nil {
		panic(err)
	}
	ModelCommand.Flags().BoolVarP(&force, "force", "f", false, "强制更新")
	ModelCommand.Flags().BoolVarP(&forceEntity, "ForceEntity", "E", false, "强制更新Entity")
	ModelCommand.Flags().BoolVarP(&forceDao, "ForceDao", "D", false, "强制更新Dao")
	ModelCommand.Flags().BoolVarP(&all, "all", "a", false, "生成全部表")
	ModelCommand.Flags().StringVarP(&upperCreateTime, "CreateTime", "C", "CreateTime", "数据表创建时间")
	ModelCommand.Flags().StringVarP(&upperUpdateTime, "UpdateTime", "U", "UpdateTime", "数据表更新时间")
}

// ModelCommand
// go run .\bin\generator.go model --table=school --conn=Default --prefix=ts_ --out=app\models
// go run .\bin\generator.go model --table=all --conn=Default --prefix=ts_ --out=app\models
var ModelCommand = &cobra.Command{
	Use:     "model",
	Short:   "Model Build Tool",
	Long:    "Model Build Tool",
	Example: "model -t=school -o=app\\models",
	Run: func(cmd *cobra.Command, args []string) {
		if table == "" && !all {
			fmt.Println("Please specify table name or all, eg: -t=table name Or --all")
			err := cmd.Help()
			if err != nil {
				panic(err)
			}
			return
		}
		var (
			isGenerateDao        YesNo = Yes
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
		if (Yes != isGenerateDao && No != isGenerateDao) && n > 0 {
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
			isGenerateEntity = Yes
			goto loopChoiceGenerateEntity
		}
	loopEntity:
		fmt.Print("Whether Generate Entity [yes|no] (default:yes):")
		n, err = fmt.Scanln(&isGenerateEntity)
		if nil != err && n > 0 {
			panic(err)
		}
		if (Yes != isGenerateEntity && No != isGenerateEntity) && n > 0 {
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
		if (Yes != confirmBeginGenerate && No != confirmBeginGenerate) && n > 0 {
			goto loopConfirmBeginGenerate
		}
		if (Yes != confirmBeginGenerate) && n > 0 {
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
		model := newModel(db, prefix, isGenerateDao, isGenerateEntity).
			setOutEntityDir(outEntityDir).
			setOutDaoDir(outDaoDir).
			setForce(force).
			setForceEntity(forceEntity).
			setForceDao(forceDao)
		if table != "" {
			model.getTable(table)
		} else if all {
			model.getTables()
		} else {
			panic("Please specify table name or all")
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
	equals              []string
	withProperty        []string
	upperProperties     []string
	newFileBuf          *bytes.Buffer
	isGenerateDao       YesNo
	isGenerateEntity    YesNo
	outEntityDir        string
	outDaoDir           string
	force               bool
	forceEntity         bool
	forceDao            bool
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
	this.equals = make([]string, 0)
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

func (this *Model) setForce(force bool) *Model {
	this.force = force
	return this
}

func (this *Model) setForceEntity(forceEntity bool) *Model {
	this.forceEntity = forceEntity
	return this
}

func (this *Model) setForceDao(forceDao bool) *Model {
	this.forceDao = forceDao
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
			upperProperty := utils.String.CaseToCamel(field.Field)
			lowerProperty := utils.String.Lcfirst(upperProperty)
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
			if propertyType == timeImport {
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
			this.mergeEqual(upperProperty)
			properties = append(properties, newProperty(isPrimaryKey, upperProperty, lowerProperty, propertyType, field.Field, field.Comment))
		}
		if fieldMaxLen > 0 {
			if !isBreak {
				isBreak = true
				goto begin
			}
		}
		this.stubContext = this.originalStubContext
		modelPackage := utils.String.CaseToCamel(strings.Replace(t.Name, this.prefix, "", 1))
		this.outfile = outDir + string(os.PathSeparator) + modelPackage + string(os.PathSeparator) + this.modelFilename
		this.replacePackage(modelPackage)
		this.replaceImport()
		this.replaceFields()
		this.replaceTableComment()
		this.replaceProperty()
		this.replaceTableName(t.Name)
		this.replaceUpperPrimaryKey(upperPrimaryKey)
		this.replaceWithProperty()
		this.replaceEquals()
		if _, err := os.Stat(this.outfile); os.IsNotExist(err) {
			this.write(this.outfile, this.stubContext)
			fmt.Println(fmt.Sprintf("Model IDE %s was created.", this.outfile))
		} else {
			if this.force {
				this.write(this.outfile, this.stubContext)
				fmt.Println(fmt.Sprintf("Model IDE %s was forced updated.", this.outfile))
			} else {
				this.oldAstEach(this.newAstEach())
				fmt.Println(fmt.Sprintf("Model IDE %s was updated.", this.outfile))
			}
		}

		entityPackage := fmt.Sprintf("%sEntity", modelPackage)

		goMod := GetModInfo()
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		childPath := GetGoModChildPath(pwd)
		var childPathStr string
		if len(childPath) > 0 {
			childPathStr = fmt.Sprintf("/%s/", strings.Join(childPath, "/"))
		}
		var goModModulePath = goMod.Module.Path
		regex := regexp.MustCompile(`/$`)
		matched := regex.MatchString(goModModulePath)
		if !matched {
			goModModulePath += "/"
		}
		modelImport := goModModulePath + fmt.Sprintf("%s%s", childPathStr, strings.ReplaceAll(utils.Dir.Dirname(this.outfile), "\\", "/"))
		modelImport = strings.ReplaceAll(modelImport, "//", "/")

		if this.isGenerateDao.Bool() {
			newEntity().
				setOutDir(this.outEntityDir).
				setPackage(entityPackage).
				setModelInfo(newModelInfo(modelImport, modelPackage)).
				setTable(t).
				setUpperPrimaryKey(upperPrimaryKey).
				setProperties(properties).
				setFieldMaxLen(fieldMaxLen).
				setPropertyTypeMaxLen(propertyTypeMaxLen).
				setUpperPropertyMaxLen(upperPropertyMaxLen).
				setForce(this.forceEntity).
				gen()

			entityImport := goModModulePath + fmt.Sprintf("%s%s", childPathStr, strings.ReplaceAll(fmt.Sprintf("%s/%s", this.outEntityDir, entityPackage), "\\", "/"))
			daoFilename := fmt.Sprintf("%sDao.go", modelPackage)
			entityImport = strings.ReplaceAll(entityImport, "//", "/")

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
				setForce(this.forceDao).
				gen()
		} else if this.isGenerateEntity.Bool() {
			newEntity().
				setOutDir(this.outEntityDir).
				setPackage(entityPackage).
				setModelInfo(newModelInfo(modelImport, modelPackage)).
				setTable(t).
				setUpperPrimaryKey(upperPrimaryKey).
				setProperties(properties).
				setFieldMaxLen(fieldMaxLen).
				setPropertyTypeMaxLen(propertyTypeMaxLen).
				setUpperPropertyMaxLen(upperPropertyMaxLen).
				setForce(this.forceEntity).
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
	if !findProperty(upperProperty, this.upperProperties) {
		this.upperProperties = append(this.upperProperties, upperProperty)
	}
}

func findProperty(upperProperty string, upperProperties []string) bool {
	for _, s := range upperProperties {
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

func (this *Model) mergeEqual(upperProperty string) {
	stub := stubs.NewStub(modelEqualStubFilename).Context()
	if len(this.equals) == 0 {
		stub = strings.ReplaceAll(stub, "%UPPER_PROPERTY%", upperProperty)
	} else {
		stub = LeftStrPad(strings.ReplaceAll(stub, "%UPPER_PROPERTY%", upperProperty), 8, " ")
	}
	this.equals = append(this.equals, stub)
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
	if token.IsKeyword(lowerProperty) {
		lowerProperty = fmt.Sprintf("_%s", lowerProperty)
	}
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
		LeftStrPad(armImport, 4, " "),
		LeftStrPad(himImport, 4, " "),
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

func (this *Model) replaceEquals() {
	this.stubContext = strings.ReplaceAll(this.stubContext, "%EQUALS%", strings.Join(this.equals, " &&\n"))
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
	gormDB := this.db.Raw(fmt.Sprintf("SHOW FULL COLUMNS FROM `%s`", tableName)).Get(&fields)
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
	return fields
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
