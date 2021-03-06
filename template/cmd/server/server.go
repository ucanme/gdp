package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ucanme/fastgo/cmd/server/middleware"
	"github.com/ucanme/fastgo/conf"
	v1 "github.com/ucanme/fastgo/controller/v1"
	"github.com/ucanme/fastgo/cron"
	"github.com/ucanme/fastgo/library/db"
	"github.com/ucanme/fastgo/library/log"
	"github.com/urfave/cli"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Server ...
var Server = cli.Command{
	Name:  "server",
	Usage: "http server",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Value: "config.toml",
			Usage: "toml配置文件",
		},
		cli.StringFlag{
			Name:  "args",
			Value: "",
			Usage: "multiconfig cmd line args",
		},
	},
	Action: run,
}

func run(c *cli.Context) {
	conf.Init(c.String("conf"), c.String("args"))
	db.Init()
	log.Init()

	srv := &http.Server{
		Handler:      GetEngine(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 60 * 5,
	}

	l, err := net.Listen("tcp4", conf.Config.Server.Listen)
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Println("server start,listen:", conf.Config.Server.Listen)
		err := srv.Serve(l)
		if err != nil {
			panic(err)
			os.Exit(1)
		}
	}()

	cron.Init()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	select {
	case s := <-sig:
		time.Sleep(time.Second) //wait for all Close function return, but not for sure
		fmt.Printf("catch signal [%s], exit now...\n", s)
	}
	if err != nil {
		panic(err)
	}
}

func GetEngine() *gin.Engine {
	r := gin.Default()

	var gRoute gin.IRouter

	gRoute = r
	gRoute.Use(middleware.Recovery())
	gRoute.Use(middleware.Access())
	if conf.Config.Auth.Enable {
		gRoute.Use(middleware.Auth())
	}

	gRoute.Use(func(c *gin.Context) {
		c.Next()
	})
	fmt.Println("hello world")
	V1(r)
	gRoute.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	return r
}

func V1(r gin.IRouter) {
	g := r.Group("/v1")
	{
		g.POST("/demo",v1.Demo)
	}
}
