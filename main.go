package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"gitee.com/ucanme/gdp/gdp"
	"github.com/urfave/cli"
	_ "embed"
)
//go:embed template statics
var Files embed.FS


func main() {
	app := cli.NewApp()
	app.Version = "gbp:0.0.1"
	app.Usage = "Generate default project layout for Go."
	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   " Generate go project default layout",
			Action: func(c *cli.Context) error {
				gdpApp,err := gdp.New(false,Files)
				if err!=nil{
					fmt.Println(err)
					return err
				}
				gdpApp.Generate(Files)
				return err
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
