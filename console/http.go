package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpFetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("httpFetch: get url=%v: %v", url, err)
	}
	defer resp.Body.Close()

	var info []byte
	info, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("httpFetch: read all: url=%v: %v", url, err)
	}

	return info, nil
}
