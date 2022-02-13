package api

import (
	"encoding/json"
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/client"
	"github.com/sav4enk0r0man/stolon-consul-discovery/logger"
	"net/http"
)

type Service struct {
	ID      string   `json:"ID"`
	Name    string   `json:"Name"`
	Address string   `json:"Address"`
	Tags    []string `json:"Tags"`
}

func RegisterService(node Node, serviceName string, conf map[string]string) (resp *http.Response, err error) {
	url := conf["url"]
	namePattern := conf["namepattern"]
	api := fmt.Sprintf("%s/v1/agent/service/register", url)

	var sevice = Service{
		ID:      fmt.Sprintf("%s-%s", node.Name(), serviceName),
		Name:    fmt.Sprintf(namePattern, serviceName),
		Address: node.Address(),
		Tags: []string{
			node.Role(),
		},
	}

	body, err := json.Marshal(sevice)
	if err != nil {
		return nil, logger.Wrapper(err, err.Error())
	}

	resp, err = client.Put(api, body, map[string]string{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DeregisterService(node Node, serviceName string, conf map[string]string) (resp *http.Response, err error) {
	url := conf["url"]
	api := fmt.Sprintf("%s/v1/agent/service/deregister/%s-%s", url, node.Name(), serviceName)

	resp, err = client.Put(api, []byte{}, map[string]string{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
