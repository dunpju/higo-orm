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
		childPath = append(childPath, utils.Dir.Basename(targetPath))
		targetPath = abovePath
		goto begin
	}
	return childPath
}
