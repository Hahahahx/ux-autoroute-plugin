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

// 校验是否是路由，并返回对应匹配到的字符串，规则很简单 a-z的字符串，后续可以跟冒号字段，或者只有冒号字段
func isRoute(name string) (string, bool) {
	r, _ := regexp.Compile("([a-z]+)?(:[a-z]+)?")
	res := r.FindStringSubmatch(name)

	if len(res) > 0 {
		return res[0], true
	} else {
		return "", false
	}
}

// 校验是否是默认路由，规则以感叹号起始
func isDefault(name string) bool {
	r, _ := regexp.Compile("^!")
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
