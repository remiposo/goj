package client

import (
	"fmt"
	"html"
	"strings"

	"github.com/remiposo/goj/model"
)

type Client interface {
	FetchSamples(contest, task string) ([]*model.Sample, error)
	ListTasks(contest string) ([]string, error)
}

func fetchSamplesErr(msg string) error {
	return fmt.Errorf("failed to fetch samples: %s", msg)
}

func listTasksErr(msg string) error {
	return fmt.Errorf("failed to get task names: %s", msg)
}

func normalizeHTMLText(rawText string) string {
	unescaped := html.UnescapeString(rawText)
	return strings.ReplaceAll(unescaped, "\r\n", "\n")
}
