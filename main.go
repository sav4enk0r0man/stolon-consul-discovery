package main

import (
	"encoding/json"
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/api"
	"log"
	"os"
)

const (
	url         = "http://127.0.0.1:8500"
	clusterName = "test2-stolon-cluster"
	index       = 0
)

var InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	InfoLog.Println("Starting...")
	response := api.WaintIndex(clusterName, url, index)
	indexes := make([]api.Index, 0)
	err := json.Unmarshal(response, &indexes)
	if err != nil {
		ErrorLog.Printf("unmarshal error: %s", err)
	}
	fmt.Printf("raw: %s", response)
	fmt.Printf("decoded: %#v", indexes[0])
}
