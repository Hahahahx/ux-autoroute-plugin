package parse

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/fatih/color"
)

type Router struct {
	Component    string `json:"element"`
	Handles      string `json:"handles"`
	RelativePath string `json:"relative"`
	AbsolutePath string `json:"absolute"`
	Path         string `json:"path"`
	PathName     string `json:"pathName"`
	Index        bool   `json:"index"`
	Lazy         bool   `json:"lazy"`
	Param        bool   `json:"param"`
	Recursion    bool   `json:"recursion"`

	// 内部使用
	relative string
	absolute string
	path     string
	pathName string

	Child []Router `json:"child"`
}

func (r *Router) ToJavaScriptString() {
	r.RelativePath = ToJavaScriptString(r.relative)
	r.AbsolutePath = ToJavaScriptString(r.absolute)
	r.Path = ToJavaScriptString(r.path)
	r.PathName = ToJavaScriptString(r.pathName)
}

func Parse(
	output,
	recursion,
	lazyImport string,
) ([]Router, []string) {

	var router Router
	router.Recursion = true
	// windows下文件路径可能会出现\\将其替换，统一为/
	router.absolute = recursion
	router.relative, _ = GetRelativePath(output, recursion)
	router.pathName = ""
	router.path = ""
	router.ToJavaScriptString()
	router.Path = ToJavaScriptString("/")

	router.Component = "Page"

	var Import = []string{"/* eslint-disable */", "// @ts-nocheck", lazyImport}
	Import = append(Import, ImportComponent(router.Component, router.RelativePath))
	router.Handles = router.Component + "Handles.handles"
	Import = append(Import, ImportComponentHandles(router.Component+"Handles", router.RelativePath))

	// 解析路由
	router.Child, Import = RecursionFile(router, Import)

	Import = append(Import, "\n\n")

	routers := []Router{router}
	return routers, Import
}

// 获取指定目录下的所有文件和目录
// AbsolutePath 绝对路径，即文件夹位置
// RelativePath 相对路径，相对于文件输出目录的路径，即router.ts到pages的路径
// Father 父级路由
// Import 导入的文件路径
func RecursionFile(Father Router, Import []string) ([]Router, []string) {
	files, err := ioutil.ReadDir(Father.absolute)

	HandleError(err, "文件夹："+Father.absolute+"打开失败")

	var (
		router Router
		childs []Router
	)

	for _, fi := range files {

		filename := path.Base(fi.Name())

		//获取文件后缀
		ext := path.Ext(filename)

		if !fi.IsDir() && !isRouteFile(ext) {
			continue
		}
		base := strings.TrimSuffix(filename, ext)

		//获取文件名
		name := base
		name, router.Lazy = isLazy(name)
		name, router.Param = isParam(name)
		name, router.Index = isIndex(name)
		name, result := isRoute(name)

		// 如果是index文件，那么其不做当前的路由，而作为父级文件夹的路由组件
		if name == "index" {
			result = false
		}

		if result {
			router.Component = Father.Component + FirstUpper(name)
			router.Recursion = fi.IsDir()
			if router.Param {
				name = ":" + name
			}

			router.pathName = name
			router.path = Father.path + "/" + name

			// 相对路径可以不加后缀，即省略tsx,jsx等
			router.relative = Father.relative + "/" + base
			// 绝对路径使用全名称
			router.absolute = Father.absolute + "/" + filename
			router.ToJavaScriptString()
			if router.Lazy {
				Import = append(Import, ImportLazyComponent(router.Component, router.RelativePath))
			} else {
				Import = append(Import, ImportComponent(router.Component, router.RelativePath))
			}

			router.Handles = router.Component + "Handles.handles"
			Import = append(Import, ImportComponentHandles(router.Component+"Handles", router.RelativePath))

			childs = append(childs, router)
		}
	}

	for index, child := range childs {
		if child.Recursion {
			childs[index].Child, Import = RecursionFile(child, Import)
		}
	}

	return childs, Import
}

// 处理错误
func HandleError(err error, condi interface{}) {
	if err != nil {
		fmt.Println(BackGroundString(color.BgHiRed, " ERROR "), condi)
		panic(1)
	}
}
