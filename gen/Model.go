package gen

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	table  string
	conn   string
	prefix string
	out    string
)

// go run .\generator.go model --table=all --conn=Default
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
	},
}

func initModel() {
	model.Flags().StringVarP(&table, "table", "t", "", "表名,all生成所有表模型")
	err := model.MarkFlagRequired("table")
	if err != nil {
		panic(err)
	}
	model.Flags().StringVarP(&conn, "conn", "c", "", "数据库连接,默认值:Default")
	model.Flags().StringVarP(&prefix, "prefix", "p", "", "数据库前缀,如:fm_")
	model.Flags().StringVarP(&out, "out", "o", "", "模型生成目录,如:app\\models")
	generator.AddCommand(model)
}
