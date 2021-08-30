package parse

import (
	"regexp"
	"strings"
)

// 将文件夹的名称转换为标题，后续作为组件引入
// 如文件夹/login 		--->   import Login
//         /login/user	--->   import LoginUser
func titleCase(name string) string {

	if name == "" {
		return name
	}

	return convertHump(name, "/")
}

// 校验是否是路由，规则很简单，只需要a-z的组成即可
func isRoute(name string) bool {
	r, _ := regexp.Compile("^[a-z]+$")
	return r.MatchString(name)
}

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
