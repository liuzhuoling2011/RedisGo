package main

import (
	"log"
	"os"
	"redisgo/utils"
	"redisgo/web"

	"github.com/urfave/cli"
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
	app.Name = utils.GetName()
	app.Version = utils.GetVersion()
	liuzhuoling := cli.Author{
		Name:  utils.GetAuthor(),
		Email: utils.GetEmail(),
	}
	app.Authors = []cli.Author{liuzhuoling}
	app.Description = utils.GetDescription()

	app.Action = func(c *cli.Context) {
		log.Println("打开浏览器, 输入 http://localhost:51299 查看效果\n")
		_ = web.StartServer(51299, true)
	}
	app.Commands = []cli.Command{
		{
			Name:     "web",
			Usage:    "启用 web 服务",
			Category: "其他",
			Action: func(c *cli.Context) error {
				log.Printf("打开浏览器, 输入: http://localhost:%d 查看效果\n", c.Uint("port"))
				log.Println(web.StartServer(c.Uint("port"), c.Bool("access")))
				return nil
			},
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "port",
					Usage: "自定义端口",
					Value: 51299,
				},
				cli.BoolFlag{
					Name:   "access",
					Usage:  "是否允许外网访问",
					Hidden: false,
				},
			},
		},
	}
	_ = app.Run(os.Args)
}
