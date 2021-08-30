package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"parse/parse"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

func main() {

	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "router-parse",
		Version:              "v1.0.0",
		Compiled:             time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Uxxx",
				Email: "1219654535@qq.com",
			},
		},
		Copyright: "(c) 2021 ux",
		Commands: []*cli.Command{
			&cli.Command{
				Name:  "parse",
				Usage: "parse filepath to router",
				Action: func(c *cli.Context) error {

					output := c.String("o")
					output = strings.ReplaceAll(output, "\\", "/")
					parse.PathExistsIsDir(output)
					recursion := c.String("r")
					recursion = strings.ReplaceAll(recursion, "\\", "/")
					parse.PathExistsIsDir(recursion)
					filename := c.String("n")

					f, err := ioutil.ReadFile(filepath.Join(output, filename))

					router := parse.RecursionFile(output, recursion, "")

					var routers []parse.Router

					routeJson, err := json.Marshal(append(routers, router))
					parse.HandleError(err, "解析JSON出错")

					parse.ImportRoute = append(parse.ImportRoute, "import loadable from '@loadable/component';\n")

					var out bytes.Buffer
					for _, importString := range parse.ImportRoute {
						// _, err = fmt.Fprintln(file, importString)
						_, err := out.WriteString(importString + "\n")
						parse.HandleError(err, "写入数据出错:"+importString)
					}

					out.WriteString("const router=")

					err = json.Indent(&out, routeJson, "", "\t")

					parse.HandleError(err, "格式化JSON出错")

					out.WriteString("\n\nexport default router")

					finalRes := strings.ReplaceAll(out.String(), "\"", "")

					parse.HandleError(err, "读取信息失败")

					if finalRes != string(f) {
						file, err := os.OpenFile(filepath.Join(output, filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
						parse.HandleError(err, "创建文件"+filename+"出错")
						file.WriteString(finalRes)
						defer file.Close()
						fmt.Println(parse.BackGroundString(color.BgHiGreen, " SUCCED "), "路由解析完成,生成文件"+file.Name())
					} else {
						fmt.Println(parse.BackGroundString(color.BgBlue, " NONE "), "文件未发生任何改变")
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Required: true,
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "router output path",
					},
					&cli.StringFlag{
						Required: true,
						Name:     "recursion",
						Aliases:  []string{"r"},
						Usage:    "recursion path",
					},
					&cli.StringFlag{
						Name:    "filename",
						Value:   "router.js",
						Aliases: []string{"n"},
						Usage:   "recursion path",
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
