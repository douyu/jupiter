package juno

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/douyu/jupiter/pkg"
	"github.com/go-resty/resty/v2"
)

// dataSource file provider.
type dataSource struct {
	path string
}

// NewDataSource returns new dataSource.
// path: juno://ip:port/{configPath}
func NewDataSource(path string, watch bool) *dataSource {
	return &dataSource{path: path}
}

type configData struct {
	Content  string `json:"content"`
	Metadata struct {
		Timestamp int      `json:"timestamp"`
		Version   string   `json:"version"`
		Format    string   `json:"format"`
		Paths     []string `json:"paths"`
	} `json:"metadata"`
}

func (ds *dataSource) ReadConfig() (content []byte, err error) {
	url, err := ds.genURL()
	if err != nil {
		return nil, err
	}

	res, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}

	return res.Body(), nil
}

func (ds *dataSource) genURL() (string, error) {
	urlParse, err := url.Parse(ds.path)
	if err != nil {
		return "", err
	}

	agent := urlParse.Host
	if agent == "" {
		return "", errors.New("juno agent adress is empty")
	}
	env := urlParse.Query().Get("env")
	if env == "" {
		return "", errors.New("juno env is empty")
	}

	return fmt.Sprintf("http://%s/api/config/%s?name=%s&env=%s",
		agent, filepath.Base(urlParse.Path), pkg.Name(), env), nil
}

// IsConfigChanged ...
func (ds *dataSource) IsConfigChanged() <-chan struct{} {
	return nil
}

// Close ...
func (ds *dataSource) Close() error {
	return nil
}
