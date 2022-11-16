package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"github.com/remiposo/goj/model"
)

var _ Client = &Atcoder{}

type Atcoder struct {
	client *http.Client
}

func NewAtcoder() *Atcoder {
	return &Atcoder{
		client: &http.Client{},
	}
}

func (a *Atcoder) FetchSamples(contest, task string) ([]*model.Sample, error) {
	u, _ := url.Parse("https://atcoder.jp")
	u.Path = path.Join(u.Path, "contests", contest, "tasks", fmt.Sprintf("%s_%s", contest, task))
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

	inputRg := regexp.MustCompile(`(?s)<h3>Sample Input.*?</h3><pre>(.*?)</pre>`)
	outputRg := regexp.MustCompile(`(?s)<h3>Sample Output.*?</h3><pre>(.*?)</pre>`)
	inputs := inputRg.FindAllSubmatch(body, -1)
	outputs := outputRg.FindAllSubmatch(body, -1)
	if len(inputs) == 0 || len(inputs) != len(outputs) {
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
