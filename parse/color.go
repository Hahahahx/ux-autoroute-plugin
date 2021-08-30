package parse

import "github.com/fatih/color"

func BackGroundString(cBg color.Attribute, format string, a ...interface{}) string {
	c := color.New(color.FgCyan).Add(cBg)
	return c.Sprint(color.BlackString(format, a...))
}
