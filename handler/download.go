package handler

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/remiposo/goj/client"
	"github.com/urfave/cli/v2"
)

func downloadErr(msg string) error {
	return fmt.Errorf("failed to download sample: %s", msg)
}

func Download(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 3 {
		cli.ShowSubcommandHelp(cCtx)
		fmt.Fprintln(cCtx.App.ErrWriter, "")
		return downloadErr("invalid arg length")
	}
	service := cCtx.Args().Get(0)
	contest := cCtx.Args().Get(1)
	task := cCtx.Args().Get(2)

	rootDir, err := findRoot(".")
	if err != nil {
		return downloadErr(err.Error())
	}

	var c client.Client
	switch service {
	case "atcoder":
		c = client.NewAtcoder()
	default:
		return downloadErr("invalid service name")
	}

	samples, err := c.FetchSamples(contest, task)
	if err != nil {
		return downloadErr(err.Error())
	}
	testDir := filepath.Join(rootDir, service, contest, task, "test")
	if _, err := os.Stat(testDir); err == nil {
		return downloadErr(fmt.Sprintf("'%s' already exists", testDir))
	}
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return downloadErr(fmt.Sprintf("failed to create '%s'", testDir))
	}

	for idx, sample := range samples {
		inputPath := filepath.Join(testDir, fmt.Sprintf("sample-%v.input", idx))
		if err := os.WriteFile(inputPath, []byte(sample.Input), 0644); err != nil {
			return downloadErr(fmt.Sprintf("failed to create '%s'", inputPath))
		}
		outputPath := filepath.Join(testDir, fmt.Sprintf("sample-%v.output", idx))
		if err := os.WriteFile(outputPath, []byte(sample.Output), 0644); err != nil {
			return downloadErr(fmt.Sprintf("failed to create '%s'", inputPath))
		}
	}
	fmt.Fprintf(cCtx.App.Writer, "downloaded samples in '%s'\n", testDir)

	return nil
}
