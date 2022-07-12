package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("code:%d,body:%s", response.StatusCode, string(body))
	}
	return body, nil
}
