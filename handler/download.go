package handler

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/remiposo/goj/client"
	"github.com/remiposo/goj/model"
	"github.com/urfave/cli/v2"
)

func downloadErr(err error) error {
	return fmt.Errorf("failed to download sample: %w", err)
}

func (h *Handler) Download(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 1 {
		cli.ShowSubcommandHelp(cCtx)
		fmt.Fprintln(cCtx.App.ErrWriter, "")
		return downloadErr(errors.New("invalid arg length"))
	}

	href, err := url.Parse(cCtx.Args().First())
	if err != nil {
		return downloadErr(err)
	}

	var samples []*model.Sample
	switch href.Host {
	case client.AtcoderHost:
		h.OJ, err = client.NewAtcoder()
		samples, err = h.OJ.FetchSamples(href)
		if err != nil {
			return downloadErr(err)
		}
	default:
		return downloadErr(errors.New("invalid hostname"))
	}

	curDir, err := os.Getwd()
	if err != nil {
		return downloadErr(err)
	}
	testDir := filepath.Join(curDir, "test")
	if _, err := os.Stat(testDir); err == nil {
		return downloadErr(fmt.Errorf("'%v' already exists", testDir))
	}
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return downloadErr(err)
	}
	for idx, sample := range samples {
		inputPath := filepath.Join(testDir, fmt.Sprintf("sample-%v.input", idx))
		outputPath := filepath.Join(testDir, fmt.Sprintf("sample-%v.output", idx))
		err := os.WriteFile(inputPath, []byte(sample.Input), 0644)
		err = os.WriteFile(outputPath, []byte(sample.Output), 0644)
		if err != nil {
			return downloadErr(err)
		}
	}
	fmt.Fprintf(cCtx.App.Writer, "successfully downloaded samples to '%v'\n", testDir)

	return nil
}
