package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dunpju/higo-utils/utils"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type YesNo string

func (this YesNo) Bool() bool {
	lower := strings.ToLower(string(this))
	if lower == Yes {
		return true
	} else if lower == No {
		return false
	}
	panic(fmt.Errorf("undefined Constant"))
}

// 转换字段类型
func convertFiledType(field TableField) string {
	types := strings.Split(field.Type, "(")
	switch types[0] {
	case "int", "smallint", "tinyint", "mediumint", "year":
		return "int"
	case "bit":
		return "byte"
	case "bigint":
		return "int64"
	case "decimal", "double", "float", "real", "numeric":
		return "float32"
	case "timestamp", "time", "datetime", "date":
		return "time.Time"
	case "binary", "varbinary":
		return "[]byte"
	case "char", "varchar", "text", "longtext", "mediumtext", "set", "enum":
		return "string"
	case "boolean":
		return "bool"
	default:
		return "interface{}"
	}
}

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

// LeftStrPad
// input string 原字符串
// padLength int 规定补齐后的字符串位数
// padString string 自定义填充字符串
func LeftStrPad(input string, padLength int, padString string) string {
	output := ""
	for i := 1; i <= padLength; i++ {
		output += padString
	}
	return output + input
}

func GetModInfo() *GoMod {
	cmd := exec.Command("go", "mod", "edit", "-json")
	buffer := bytes.NewBufferString("")
	cmd.Stdout = buffer
	cmd.Stderr = buffer

	if err := cmd.Run(); err != nil {
		panic(err)
	}
	goMod := &GoMod{}
	err := json.Unmarshal(buffer.Bytes(), &goMod)
	if err != nil {
		panic(err)
	}
	return goMod
}

type GoMod struct {
	Module  Module
	Go      string
	Require []Require
	Exclude []Module
}

type Module struct {
	Path    string
	Version string
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

func Dirslice(path string) []string {
	pathSeparator := string(os.PathSeparator)
	paths := strings.Split(path, pathSeparator)
	if len(paths) == 1 {
		re := regexp.MustCompile("/")
		if re.Match([]byte(paths[0])) {
			paths = strings.Split(path, pathSeparator)
		}
	}
	return paths
}
