package handler

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/remiposo/goj/model"
	"github.com/urfave/cli/v2"
)

func testErr(err error) error {
	return fmt.Errorf("failed to test: %w", err)
}

func Test(cCtx *cli.Context) error {
	if cCtx.Args().Len() == 0 {
		cli.ShowSubcommandHelp(cCtx)
		fmt.Fprintln(cCtx.App.ErrWriter, "")
		return testErr(errors.New("invalid arg length"))
	}

	entries, err := os.ReadDir("test")
	if err != nil {
		return testErr(err)
	}
	samples := make(map[string]*model.Sample)
	for _, entry := range entries {
		baseName := entry.Name()
		r := regexp.MustCompile(`^sample-(\d+)\.(input|output)$`)
		matches := r.FindSubmatch([]byte(baseName))
		if len(matches) == 0 {
			continue
		}

		body, err := os.ReadFile(filepath.Join("test", baseName))
		if err != nil {
			return testErr(err)
		}
		sample, ok := samples[string(matches[1])]
		if !ok {
			sample = &model.Sample{}
			samples[string(matches[1])] = sample
		}
		if string(matches[2]) == "input" {
			sample.Input = string(body)
		} else {
			sample.Output = string(body)
		}
	}

	args := cCtx.Args().Slice()
	for id, sample := range samples {
		fmt.Fprintf(cCtx.App.Writer, "=== RUN: sample-%v\n", id)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = bytes.NewBuffer([]byte(sample.Input))
		got, err := cmd.Output()
		if err != nil {
			fmt.Fprintf(cCtx.App.Writer, "=== FAIL: sample-%v\n", id)
			fmt.Fprintf(cCtx.App.Writer, "Runtime Error: %v\n", err)
			continue
		}
		if sample.Output != string(got) {
			fmt.Fprintf(cCtx.App.Writer, "=== FAIL: sample-%v\n", id)
			fmt.Fprintf(cCtx.App.Writer, "want %v\ngot %v\n", sample.Output, string(got))
			continue
		}

		fmt.Fprintf(cCtx.App.Writer, "=== PASS: sample-%v\n", id)
	}

	return nil
}
