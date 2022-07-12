package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func HttpGetReturnReqAndResp(url string) ([]byte, *http.Request, *http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, request, nil, err
	}
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, request, response, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, request, response, err
	}
	return body, request, response, nil
}

func HttpGetJsonOrContentType(url string) ([]byte, string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()
	contentType := response.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, contentType, nil
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}
	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("code:%d,body:%s", response.StatusCode, string(body))
	}
	return body, "", nil
}
