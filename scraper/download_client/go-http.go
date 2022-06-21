package download_client

import (
	"io/ioutil"
	"net/http"
)

type GoHttp struct {
}

func (dClient GoHttp) Open() error {
	return nil
}

func (dClient GoHttp) Close() error {
	return nil
}

func (dClient GoHttp) Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	htmlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return htmlBytes, nil
}

func NewGoHttpDownloadClient(config map[string]interface{}) (DownloadClient, error) {
	return GoHttp{}, nil
}
