package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitee.com/ucanme/scaffold/scaffold"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0-rc"
	app.Usage = "Generate scaffold project layout for Go."
	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   " Generate scaffold project layout",
			Action: func(c *cli.Context) error {
				path, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				currDir, err := filepath.Abs(path)
				if err != nil {
					return err
				}

				err = scaffold.New(false).Generate(currDir)
				//fmt.Printf("error:%+v\n", err)
				if err == nil {
					fmt.Println("Success Created. Please excute `make up` to start service.")
				}

				return err
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
