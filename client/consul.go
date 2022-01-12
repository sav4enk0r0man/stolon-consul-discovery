package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func ConsulClient(url string) []byte {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "api http request: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "api http request: read url %s: %v\n", url, err)
		os.Exit(1)
	}

	return body
}
