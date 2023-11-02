package gen

import (
	"github.com/dunpju/higo-utils/utils"
	"path/filepath"
)

func GetGoModChildPath(targetPath string) []string {
	childPath := make([]string, 0)
begin:
	abovePath := utils.Dir.Dirname(targetPath)
	files, err := filepath.Glob(targetPath + "/go.mod")
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		path := []string{utils.Dir.Basename(targetPath)}
		childPath = append(path, childPath...)
		targetPath = abovePath
		goto begin
	}
	return childPath
}
