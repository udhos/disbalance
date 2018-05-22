package main

import (
	"bytes"
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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("httpFetch: bad status: %d", resp.StatusCode)
	}

	info, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, fmt.Errorf("httpFetch: read all: url=%v: %v", url, errRead)
	}

	return info, nil
}

func httpPost(url, contentType string, buf []byte) ([]byte, error) {
	resp, err := http.Post(url, contentType, bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("httpPost: url=%v: %v", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("httpPost: bad status: %d", resp.StatusCode)
	}

	info, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, fmt.Errorf("httpPost: read: url=%v: %v", url, errRead)
	}

	return info, nil
}

func httpPut(url, contentType string, buf []byte) ([]byte, error) {
	req, errNew := http.NewRequest("PUT", url, bytes.NewReader(buf))
	if errNew != nil {
		return nil, errNew
	}

	req.Header.Set("Content-Type", contentType)

	client := http.Client{}

	resp, errDel := client.Do(req)
	if errDel != nil {
		return nil, errDel
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("httpPut: bad status: %d", resp.StatusCode)
	}

	info, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}

	return info, nil
}

func httpDelete(url string) ([]byte, error) {
	req, errNew := http.NewRequest("DELETE", url, nil)
	if errNew != nil {
		return nil, errNew
	}

	client := http.Client{}

	resp, errDel := client.Do(req)
	if errDel != nil {
		return nil, errDel
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("httpDelete: bad status: %d", resp.StatusCode)
	}

	info, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}

	return info, nil
}
