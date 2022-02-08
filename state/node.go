package state

import (
	"github.com/sav4enk0r0man/stolon-consul-discovery/api"
	"net/http"
)

func GetMasters(s api.ClusterState) []api.Node {
	var masters []api.Node
	for _, node := range s.Nodes() {
		if node.Role() == "master" {
			masters = append(masters, node)
		}
	}

	return masters
}

func GetAddresses(s api.ClusterState) []string {
	var addresses []string
	for _, node := range s.Nodes() {
		addresses = append(addresses, node.Address())
	}
	return addresses
}

func Deregister(node api.Node, url string) *http.Response {
	return api.DeregisterService(node, url)
}

func Register(node api.Node, serviceName, url string) *http.Response {
	return api.RegisterService(node, serviceName, url)
}
