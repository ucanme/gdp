package main

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"gitee.com/ucanme/gdp/gdp"
	"github.com/urfave/cli"
	"log"
	"os"
)
//go:embed template
var Files embed.FS


func main() {
	app := cli.NewApp()
	app.Version = "gbp:0.0.1"
	app.Usage = "Generate default project layout for Go."

	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"new"},
			Usage:   " Generate go project default layout",
			Action: func(c *cli.Context) error {
				appName := c.String("name")
				if  appName == ""{
					appName = confirm("please input your app name:")
				}

				if appName == ""{
					return errors.New("the app name input not correct")
				}

				gdpApp, err := gdp.New(appName,false, Files)
				if err != nil {
					fmt.Println(err)
					return err
				}
				gdpApp.Generate(Files)
				return err
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "the app name will be generated",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}


func confirm(msg string) string {
	fmt.Printf(msg)
	var response string
	n, err := fmt.Scanln(&response)
	if err != nil {
		return ""
	}
	if n == 0{
		return ""
	}
	return response
}