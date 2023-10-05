package gen

import (
	"fmt"
	"github.com/dunpju/higo-orm/gen/stubs"
	"github.com/dunpju/higo-orm/him"
	"github.com/spf13/cobra"
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
	modelStubFilename             = "model.stub"
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
		fmt.Println(os.Getwd())
		fmt.Println(table)
		fmt.Println(conn)
		fmt.Println(prefix)
		fmt.Println(out)
		db, err := him.DBConnect(conn)
		if err != nil {
			panic(err)
		}
		if prefix == "" {
			prefix = db.DBC().Prefix()
		}
		model := newModel(db)
		model.GetTableFields(table)
		model.GetTables(prefix)
		fmt.Println(stubs.NewStub(modelStubFilename).Context())
	},
}

type Model struct {
	db     *him.DB
	fields []TableField
}

func newModel(db *him.DB) *Model {
	return &Model{db: db, fields: make([]TableField, 0)}
}

// GetTables 获取数据库表
func (this *Model) GetTables(prefix string) []Table {
	var tables []Table
	gormDB := this.db.Raw(fmt.Sprintf(`SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='%s' AND TABLE_NAME LIKE '%s%%'`, this.db.DBC().Database(), prefix)).Get(&tables)
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
	return tables
}

// GetTableFields 获取表所有字段信息
func (this *Model) GetTableFields(tableName string) []TableField {
	gormDB := this.db.Raw(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", tableName)).Get(&this.fields)
	if gormDB.Error != nil {
		panic(gormDB.Error)
	}
	return this.fields
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

// 获取字段类型
func getFiledType(field TableField) string {
	if field.Null == "YES" {
		return "interface{}"
	}
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
	default:
		return "string"
	}
}
