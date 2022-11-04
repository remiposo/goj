package main

import (
	"log"
	"os"

	"github.com/remiposo/goj/handler"
	"github.com/urfave/cli/v2"
)

func main() {
	h := &handler.Handler{}
	app := &cli.App{
		Name:  "goj",
		Usage: "simple/portable CLI tool for online judge sites",
		Commands: []*cli.Command{
			{
				Name:            "download",
				Aliases:         []string{"d"},
				Usage:           "download samples",
				Action:          h.Download,
				HideHelpCommand: true,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
