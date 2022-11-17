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

func download(cCtx *cli.Context, c client.Client, serviceDir, contest, task string) error {

	testDir := filepath.Join(serviceDir, contest, task, testDir)
	if _, err := os.Stat(testDir); err == nil {
		return downloadErr(fmt.Sprintf("'%s' already exists", testDir))
	}
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return downloadErr(fmt.Sprintf("failed to create '%s'", testDir))
	}

	samples, err := c.FetchSamples(contest, task)
	if err != nil {
		return downloadErr(err.Error())
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

func Download(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 2 && cCtx.Args().Len() != 3 {
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

	serviceDir := filepath.Join(rootDir, service)

	var c client.Client
	switch service {
	case "atcoder":
		c = client.NewAtcoder()
	default:
		return downloadErr("invalid service name")
	}

	if task != "" {
		if err := download(cCtx, c, serviceDir, contest, task); err != nil {
			return err
		}
		return nil
	}

	tasks, err := c.ListTasks(contest)
	if err != nil {
		return downloadErr(err.Error())
	}
	for _, task := range tasks {
		if err := download(cCtx, c, serviceDir, contest, task); err != nil {
			return err
		}
	}
	return nil
}
