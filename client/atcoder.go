package client

import (
	"errors"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/remiposo/goj/model"
)

const (
	AtcoderHost = "atcoder.jp"
)

var _ OJ = &Atcoder{}

type Atcoder struct {
	Client *http.Client
}

func NewAtcoder() (*Atcoder, error) {
	return &Atcoder{Client: &http.Client{}}, nil
}

func (a *Atcoder) FetchSamples(href *url.URL) ([]*model.Sample, error) {
	if href.Host != AtcoderHost {
		return nil, errors.New("invalid hostname")
	}

	req, err := http.NewRequest("GET", href.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	inputRg := regexp.MustCompile(`(?s)<h3>Sample Input.*?</h3><pre>(.*?)</pre>`)
	outputRg := regexp.MustCompile(`(?s)<h3>Sample Output.*?</h3><pre>(.*?)</pre>`)
	inputs := inputRg.FindAllSubmatch(body, -1)
	outputs := outputRg.FindAllSubmatch(body, -1)
	if len(inputs) == 0 || len(inputs) != len(outputs) {
		return nil, errors.New("samples not found")
	}
	samples := make([]*model.Sample, len(inputs))
	for idx := 0; idx < len(inputs); idx++ {
		samples[idx] = &model.Sample{
			Input:  html.UnescapeString(string(inputs[idx][1])),
			Output: html.UnescapeString(string(outputs[idx][1])),
		}
	}

	return samples, nil
}
