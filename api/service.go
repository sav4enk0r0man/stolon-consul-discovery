package api

import (
	"encoding/json"
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/client"
	"net/http"
)

type Service struct {
	ID      string   `json:"ID"`
	Name    string   `json:"Name"`
	Address string   `json:"Address"`
	Tags    []string `json:"Tags"`
}

func RegisterService(node Node, serviceName, url string) *http.Response {
	api := fmt.Sprintf("%s/v1/agent/service/register", url)

	var sevice = Service{
		ID:      node.Name(),
		Name:    fmt.Sprintf("postgresql-%s", serviceName),
		Address: node.Address(),
		Tags: []string{
			node.Role(),
		},
	}

	body, err := json.Marshal(sevice)
	if err != nil {
		logError.Fatalf("marshal error: %s", err)
	}

	resp, err := client.Put(api, body, map[string]string{})
	if err != nil {

	}
	return resp
}

func DeregisterService(node Node, url string) *http.Response {
	api := fmt.Sprintf("%s/v1/agent/service/deregister/%s", url, node.Name())
	resp, err := client.Put(api, []byte{}, map[string]string{})
	if err != nil {

	}
	return resp
}
