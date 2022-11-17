package main

import (
	"fmt"
	"os"

	"github.com/remiposo/goj/handler"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "goj",
		Usage: "simple/portable CLI tool for online judge sites",
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "create new goj project",
				Action:  handler.Create,
			},
			{
				Name:    "download",
				Aliases: []string{"d"},
				Usage:   "download samples",
				Action:  handler.Download,
			},
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "test with samples",
				Action:  handler.Test,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(app.ErrWriter, err)
		os.Exit(1)
	}
}
