package parse

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// 获取从A到B的相对路径
// 例如：from=/a/b/c/d，to=/a/e/f 相对路径为 ../../e/f
// 参数from，to必须是绝对路径
func GetRelativePath(from, to string) (string, error) {

	from, _ = filepath.Abs(from)
	to, _ = filepath.Abs(to)

	from = strings.ReplaceAll(from, "\\", "/")
	to = strings.ReplaceAll(to, "\\", "/")

	// fmt.Println(from, "==>", to)

	if from == "" || to == "" {
		return "", errors.New("绝对路径不可以为空")
	}

	arrFrom := strings.Split(from[1:], "/")
	arrTo := strings.Split(to[1:], "/")
	depth := 0
	for i := 0; i < len(arrFrom) && i < len(arrTo); i++ {
		if arrFrom[i] == arrTo[i] {
			depth++
		} else {
			break
		}
	}
	prefix := ""
	if len(arrFrom)-depth-1 <= 0 {
		prefix = "./"
	} else {
		for i := len(arrFrom) - depth - 1; i > 0; i-- {
			prefix += "../"
		}
	}
	// fmt.Println(depth)
	if len(arrTo)-depth > 0 {
		prefix += strings.Join(arrTo[depth:], "/")
	}
	return prefix, nil
}

// 检查目录是否存在
func PathExistsIsDir(path string) {

	fi, err := os.Stat(path)
	HandleError(err, "目录不存在")

	if !fi.IsDir() {
		HandleError(errors.New("err"), "不是文件夹")
	}
}

func PrintStruct(item interface{}) {
	routeJson, err := json.Marshal(item)
	HandleError(err, "解析struct失败")
	var out bytes.Buffer
	err = json.Indent(&out, routeJson, "", "\t")
	out.WriteString("\n")
	HandleError(err, "解析struct失败")
	out.WriteTo(os.Stdout)
}

func ToJavaScriptString(str string) string {
	return "'" + str + "'"
}

func ImportLazyComponent(component, path string) string {
	return "const " + component + " = lazy(() => import(" + path + "));"
}

func ImportComponent(component, path string) string {
	return "import " + component + " from " + path + ";"
}

func ImportComponentHandles(component, path string) string {
	return "import * as " + component + " from " + path + ";"
}
