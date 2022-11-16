package handler

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func createErr(msg string) error {
	return fmt.Errorf("failed to create project: %s", msg)
}

func Create(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		cli.ShowSubcommandHelp(cCtx)
		fmt.Fprintln(cCtx.App.ErrWriter, "")
		return createErr("invalid arg length")
	}

	rootDir := filepath.Clean(cCtx.Args().First())
	if _, err := os.Stat(rootDir); err == nil {
		return createErr(fmt.Sprintf("'%s' already exists", rootDir))
	}
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return createErr(fmt.Sprintf("failed to create '%s'", rootDir))
	}
	confPath := filepath.Join(rootDir, "goj.yml")
	if err := os.WriteFile(confPath, []byte{}, 0644); err != nil {
		return createErr(fmt.Sprintf("failed to create '%s'", confPath))
	}
	fmt.Fprintf(cCtx.App.Writer, "created new project in '%v'\n", rootDir)

	return nil
}
