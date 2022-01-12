package main

import (
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/client"
)

const (
	url = "http://127.0.0.1:8500/v1/kv?keys"
)

func main() {
	response := client.ConsulClient(url)
	fmt.Printf("%s", response)
}
