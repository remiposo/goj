package client

import (
	"net/url"

	"github.com/remiposo/goj/model"
)

type OJ interface {
	FetchSamples(href *url.URL) ([]*model.Sample, error)
}
