package parse

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Router struct {
	Config    interface{} `json:"config"`
	Component string      `json:"component"`
	Path      string      `json:"path"`
	Child     []Router    `json:"child"`
}

var (
	ImportRoute []string
)

//获取指定目录下的所有文件和目录
func RecursionFile(outputPath, dirPath, routePath string, lazyImport bool) Router {
	files, err := ioutil.ReadDir(dirPath)
	HandleError(err, "文件夹："+dirPath+"打开失败")

	var (
		router Router
		child  []string
	)

	for _, fi := range files {
		// 处理config文件
		handleConfigFile(fi, outputPath, dirPath, routePath, &router, lazyImport)
		handleIndexFile(fi, outputPath, dirPath, routePath, &router, lazyImport)
		if fi.IsDir() { // 目录, 递归遍历
			child = append(child, strings.ReplaceAll(filepath.Join(dirPath, fi.Name()), "\\", "/"))
		}
	}

	// PrintStruct(router)

	for _, dir := range child {
		childRouter := RecursionFile(outputPath, dir, routePath+"/"+filepath.Base(dir), lazyImport)
		router.Child = append(router.Child, childRouter)
	}

	return router
}

// 处理配置文件
func handleConfigFile(file os.FileInfo, outputPath, dirPath, routePath string, router *Router, lazyImport bool) {

	if filepath.Base(file.Name()) == "route.config" {
		// abs, err := filepath.Abs()
		// fmt.Println(abs)
		// jsonFile, err := os.Open(filepath.Join(dirPath, file.Name()))
		// HandleError(err, "读取JSON文件失败")
		// defer jsonFile.Close()
		byteValue, err := ioutil.ReadFile(filepath.Join(dirPath, file.Name()))
		HandleError(err, "读取JSON失败")
		// fmt.Println(string(byteValue))
		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)
		router.Config = result

		// 默认配置按需加载
		if lazyImport {
			noLazy, ok := result["noLazy"]

			// 判断是否存在noLazy字段
			if ok {
				// 判断是否是bool类型
				noLazy, ok := noLazy.(bool)
				if !ok {
					HandleError(errors.New("err"), "错误的字段类型，noLazy必须bool类型")
				}
				// 静态导入组件
				if noLazy {
					reletivePath, err := getRelativePath(outputPath, dirPath)
					HandleError(err, "获取文件"+dirPath+"相对路径失败")
					router.Component = "Page" + titleCase(routePath)
					ImportRoute = append(ImportRoute, "import Page"+titleCase(routePath)+" from '"+reletivePath+"/index';")
				}
			}
		} else {
			lazy, ok := result["lazy"]
			// 判断是否存在lazy字段
			if ok {
				// 判断是否是bool类型
				lazy, ok := lazy.(bool)
				if !ok {
					HandleError(errors.New("err"), "错误的字段类型，lazy必须bool类型")
				}
				// 如果是需要按需加载则交给后面的流程来处理
				if lazy {
					return
				}
			}
			reletivePath, err := getRelativePath(outputPath, dirPath)
			HandleError(err, "获取文件"+dirPath+"相对路径失败")
			router.Component = "Page" + titleCase(routePath)
			ImportRoute = append(ImportRoute, "import Page"+titleCase(routePath)+" from '"+reletivePath+"/index';")
		}
		HandleError(err, "读取文件内容失败")
	}
}

// 处理组件
func handleIndexFile(file os.FileInfo, outputPath, dirPath, routePath string, router *Router, lazyImport bool) {

	// fullName := filepath.Base(dirPath)
	// extionName := filepath.Ext(fullName
	// clearName := strings.TrimSuffix(fullName, extionName)

	if file.Name() == "index.jsx" || file.Name() == "index.tsx" {
		if router.Component == "" {

			if lazyImport {

				reletivePath, err := getRelativePath(outputPath, dirPath)
				HandleError(err, "获取文件"+dirPath+"相对路径失败")
				router.Component = "loadable(function(){return import('" + reletivePath + "')})"
			} else {

				reletivePath, err := getRelativePath(outputPath, dirPath)
				HandleError(err, "获取文件"+dirPath+"相对路径失败")
				router.Component = "Page" + titleCase(routePath)
				ImportRoute = append(ImportRoute, "import Page"+titleCase(routePath)+" from '"+reletivePath+"/index';")
			}
		}

		router.Path = "'" + routePath + "'"
	}
}

// 处理错误
func HandleError(err error, condi interface{}) {
	if err != nil {
		fmt.Println(BackGroundString(color.BgHiRed, " ERROR "), condi)
		panic(1)
	}
}
