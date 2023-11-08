package gen

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var ModelGenerator = &cobra.Command{
	Use:   "",
	Short: "Model 构建工具",
	Long:  `Model 构建工具`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("build tools")
	},
}

func Execute() {
	InitModel()
	ModelGenerator.AddCommand(ModelCommand)
	if err := ModelGenerator.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
