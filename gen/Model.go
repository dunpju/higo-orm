package gen

import (
	"fmt"
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

func initModel() {
	model.Flags().StringVarP(&table, "table", "t", "", "表名,all生成所有表模型")
	err := model.MarkFlagRequired("table")
	if err != nil {
		panic(err)
	}
	model.Flags().StringVarP(&conn, "conn", "c", "Default", "数据库连接,默认值:Default")
	model.Flags().StringVarP(&prefix, "prefix", "p", "", "数据库前缀,如:fm_")
	model.Flags().StringVarP(&out, "out", "o", "", "模型生成目录,如:app\\models")
	generator.AddCommand(model)
}

// go run .\bin\generator.go model --table=all --conn=Default
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
		newModel(db).GetTableFields(table)
	},
}

type Model struct {
	db *him.DB
}

func newModel(db *him.DB) *Model {
	return &Model{db: db}
}

// GetTableFields 获取表所有字段信息
func (this *Model) GetTableFields(tableName string) []TableField {
	var fields []TableField
	gormDB := this.db.Raw(fmt.Sprintf("SHOW FULL COLUMNS FROM ?"), tableName).Get(&fields)
	if gormDB.Error != nil {
		panic(gormDB.Error.Error())
	}
	fmt.Println(fields)
	return fields
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
