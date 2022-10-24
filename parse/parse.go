package parse

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Router struct {
	Component     string `json:"element"`
	RelativePath  string `json:relativepath`
	RealPath      string `json:realpath`
	ComponentPath string `json:componentpath`
	Path          string `json:"path"`
	PathName      string `json:"pathName"`
	Index         bool   `json:"index"`
	Lazy          bool   `json:lazy`
	Param         bool   `json:param`
	Recursion     bool   `json:recursion`

	Child []Router `json:"child"`
}

// 获取指定目录下的所有文件和目录
// AbsolutePath 绝对路径，即文件夹位置
// RelativePath 相对路径，相对于文件输出目录的路径，即router.ts到pages的路径
// Father 父级路由
// Import 导入的文件路径
func RecursionFile(Father Router, Import []string) []Router {
	files, err := ioutil.ReadDir(Father.RealPath)

	HandleError(err, "文件夹："+Father.RealPath+"打开失败")

	var (
		router Router
		childs []Router
	)

	for _, fi := range files {

		filename := path.Base(fi.Name())
		//获取文件后缀
		ext := path.Ext(filename)

		if !isRouteFile(ext) {
			continue
		}

		//获取文件名
		name := strings.TrimSuffix(filename, ext)
		name, router.Lazy = isLazy(name)
		name, router.Param = isParam(name)
		name, router.Index = isIndex(name)
		name, result := isRoute(name)

		// 如果是index文件，那么其不做当前的路由，而作为父级文件夹的路由组件
		if name == "index" {
			result = false
		}

		if result {
			router.Recursion = fi.IsDir()
			// windows下文件路径可能会出现\\将其替换，统一为/
			router.RealPath = filepath.Join(Father.RealPath, fi.Name())
			router.PathName = name
			router.Path = Father.Path + "/" + name
			router.RelativePath = filepath.Join(Father.RelativePath, name)

			if router.Lazy {
				router.Component = "lazy(() => import('" + router.RelativePath + "'))"
			} else {
				router.Component = "Page" + titleCase(router.Path)
				Import = append(Import, "import "+router.Component+" from "+router.RelativePath)
			}

			childs = append(childs, router)
		}
	}

	for _, child := range childs {
		if child.Recursion {
			child.Child = RecursionFile(child, Import)
		}
	}

	return childs
}

// 处理错误
func HandleError(err error, condi interface{}) {
	if err != nil {
		fmt.Println(BackGroundString(color.BgHiRed, " ERROR "), condi)
		panic(1)
	}
}
