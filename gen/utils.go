package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	hasGoMod  bool
	goModPath string
)

func seekGoMod(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		fmt.Println("路径：", path)
		if info.Name() == "go.mod" {
			hasGoMod = true
			goModPath = path
		}
		fmt.Println("是否为目录：", info.IsDir())
		fmt.Println("文件名：", info.Name())
		fmt.Println("大小：", info.Size())
		fmt.Println("权限：", info.Mode())
		fmt.Println("修改时间：", info.ModTime())
		fmt.Println()
	}
	return err
}

func GetGoModPath(targetPath string) {
	err := filepath.Walk(targetPath, seekGoMod)
	if err != nil {
		panic(err)
	}
}
