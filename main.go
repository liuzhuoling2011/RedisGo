package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"redisgo/utils"
	"redisgo/web"
)

var (
	Version = "1.0.0"
)

func init() {
	if !utils.InitConfig() {
		os.Exit(-1)
	}
	if !utils.InitContainers() {
		os.Exit(-1)
	}
}

func main() {
	defer utils.SaveConfig()

	app := cli.NewApp()
	app.Name = "RedisGo"
	app.Version = Version
	liuzhuoling := cli.Author{
		Name:  "liuzhuoling",
		Email: "liuzhuoling2011@hotmail.com",
	}
	app.Authors = []cli.Author{liuzhuoling}
	app.Description = "这个软件可以让你更好的管理/监控Redis"

	app.Action = func(c *cli.Context) {
		fmt.Printf("打开浏览器, 输入 http://localhost:51299 查看效果\n")
		web.StartServer(51299, true)
	}
	app.Commands = []cli.Command{
		{
			Name:     "web",
			Usage:    "启用 web 服务",
			Category: "其他",
			Action: func(c *cli.Context) error {
				fmt.Printf("打开浏览器, 输入: http://localhost:%d 查看效果\n", c.Uint("port"))
				fmt.Println(web.StartServer(c.Uint("port"), c.Bool("access")))
				return nil
			},
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "port",
					Usage: "自定义端口",
					Value: 51299,
				},
				cli.BoolFlag{
					Name:  "access",
					Usage: "是否允许外网访问",
					Hidden: false,
				},
			},
		},
	}
	app.Run(os.Args)
}
