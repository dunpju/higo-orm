package gen

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var generator = &cobra.Command{
	Use:   "",
	Short: "构建工具",
	Long:  `构建工具`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("build tools")
	},
}

func init() {
	initModel()
}

func Execute() {
	if err := generator.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
