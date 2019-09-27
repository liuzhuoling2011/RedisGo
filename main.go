package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"redisgo/web"
)

var (
	Version = "3.6.8"
)

func main() {
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
	app.Run(os.Args)
}
