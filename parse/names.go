package parse

import (
	"regexp"
	"strings"
)

// 将文件夹的名称转换为标题，后续作为组件引入
// 如文件夹/login 		--->   import Login
//
//	/login/user	--->   import LoginUser
func titleCase(name string) string {

	if name == "" {
		return name
	}

	return convertHump(name, "/")
}

// 判断文件后缀是否是有效的路由文件，例如：oth.jsx,about.tsx,[id].mdx
func isRouteFile(ext string) bool {
	words := [...]string{".tsx", ".jsx", ".mdx", ".md"}
	for _, word := range words {
		if ext == word {
			return true
		}
	}
	return false
}

// 校验是否是路由，并返回对应匹配到的字符串，规则很简单 a-z的字符串
var isRoute = verifyBase(`([a-z]+)?`)

// 校验是否是参数路由，即文件名为[name] 后续会转换成 :name，供给前端路由组件去匹配
var isParam = verifyBase(`^\[(.*)\]$`)

// 校验是否是懒加载路由，以~开头
var isLazy = verifyBase(`^~(.*)`)

// 校验是否是默认路由，以!结尾，例如：test!.tsx
var isIndex = verifyBase(`(.*)!$`)

// ascii 大小写转换
func convertCase(c rune, t int) rune {
	if t == 0 {
		return c - 32
	} else {
		return c + 32
	}

}

// 转换为驼峰，根据flag='/'  /main/user  ->  MainUser
func convertHump(str string, flag string) string {
	for {
		index := strings.Index(str, flag)

		if index == -1 {
			break
		}

		replaceCount := len(flag)

		str = str[:index] + strings.ToUpper(string(str[index+1])) + str[index+replaceCount+1:]

	}

	return str
}

func verifyBase(exg string) func(name string) (string, bool) {
	return func(name string) (string, bool) {
		r := regexp.MustCompile(exg)
		res := r.FindStringSubmatch(name)
		if len(res) > 0 {
			return res[1], true
		} else {
			return name, false
		}
	}
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
