package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/remiposo/goj/model"
)

var _ Client = &Atcoder{}

const atcoderOrigin = "https://atcoder.jp"

type Atcoder struct {
	origin string
	client *http.Client
}

func NewAtcoder() *Atcoder {
	return &Atcoder{
		origin: atcoderOrigin,
		client: &http.Client{},
	}
}

func (a *Atcoder) taskPath(contest, task string) string {
	return strings.ReplaceAll(contest, "-", "_") + "_" + task
}

func (a *Atcoder) FetchSamples(contest, task string) ([]*model.Sample, error) {
	u, _ := url.Parse(a.origin)
	u.Path = path.Join(u.Path, "contests", contest, "tasks", a.taskPath(contest, task))
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fetchSamplesErr("invalid request")
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fetchSamplesErr("failed to get response")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fetchSamplesErr("failed to read response")
	}

	inputRg := regexp.MustCompile(`(?s)<h3[^>]*>(?:Sample Input|入力例)[^<]*</h3>\s*<pre[^>]*>([^<]*)</pre>`)
	outputRg := regexp.MustCompile(`(?s)<h3[^>]*>(?:Sample Output|出力例)[^<]*</h3>\s*<pre[^>]*>([^<]*)</pre>`)
	inputs := inputRg.FindAllSubmatch(body, -1)
	outputs := outputRg.FindAllSubmatch(body, -1)
	if len(inputs) == 0 || len(inputs) != len(outputs) {
		fmt.Println(len(inputs), len(outputs))
		fmt.Println(string(body))
		return nil, fetchSamplesErr("sample not found")
	}

	samples := make([]*model.Sample, len(inputs))
	for idx := 0; idx < len(inputs); idx++ {
		samples[idx] = &model.Sample{
			Input:  normalizeHTMLText(string(inputs[idx][1])),
			Output: normalizeHTMLText(string(outputs[idx][1])),
		}
	}

	return samples, nil
}

func (a *Atcoder) ListTasks(contest string) ([]string, error) {
	u, _ := url.Parse(a.origin)
	u.Path = path.Join(u.Path, "contests", contest, "tasks")
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, listTasksErr("invalid request")
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, listTasksErr("failed to get response")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, listTasksErr("failed to read response")
	}
	taskPath := a.taskPath(contest, `([^"]*)`)
	rg := regexp.MustCompile(fmt.Sprintf(`(?s)<a\s*href="/contests/%s/tasks/%s"\s*>`, contest, taskPath))
	ms := rg.FindAllSubmatch(body, -1)
	if len(ms) == 0 {
		return nil, listTasksErr("task not found")
	}
	tasks := make([]string, 0, len(ms))
	taskMap := make(map[string]bool)
	for _, m := range ms {
		task := string(m[1])
		if !taskMap[task] {
			taskMap[task] = true
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}
