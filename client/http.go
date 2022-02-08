package client

import (
	"bytes"
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/logger"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultHTTPTimeout = 10
)

var logError = logger.DefaultLog.Error

var httpClient *http.Client

func Get(url string, opts map[string]string) ([]byte, error) {
	if err := initClient(opts); err != nil {
		return nil, err
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, logger.Wrapper(err, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("%d", resp.StatusCode)
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, logger.Wrapper(err, err.Error())
	}

	return response, nil
}

func Put(url string, body []byte, opts map[string]string) (*http.Response, error) {
	if err := initClient(opts); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		logError.Fatalf("create http request: %v\n", err)
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	response, err := httpClient.Do(request)
	if err != nil {
		logError.Fatalf("api http request: %v\n", err)
	}
	defer response.Body.Close()

	return response, nil
}

func initClient(opts map[string]string) error {
	if httpClient == nil {
		timeout := defaultHTTPTimeout
		if opts["httptimeout"] != "" {
			var err error
			timeout, err = strconv.Atoi(opts["httptimeout"])
			if err != nil {
				return logger.Wrapper(err, err.Error())
			}
		}
		httpClient = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}
	return nil
}
