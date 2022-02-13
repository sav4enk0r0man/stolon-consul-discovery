package state

import (
	"github.com/sav4enk0r0man/stolon-consul-discovery/api"
	"strings"
)

func getServices(services string) []string {
	s := strings.Split(services, ",")
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
	}
	return s
}

func getMasters(s api.ClusterState) []api.Node {
	var masters []api.Node
	for _, node := range s.Nodes() {
		if node.Role() == "master" {
			masters = append(masters, node)
		}
	}

	return masters
}

func getAddresses(s api.ClusterState) []string {
	var addresses []string
	for _, node := range s.Nodes() {
		addresses = append(addresses, node.Address())
	}
	return addresses
}

func DeregisterServices(node api.Node, context Context) (status []string, err error) {
	conf := context.Config

	for _, service := range getServices(conf["services"]) {
		resp, err := api.DeregisterService(node, service, conf)
		if err != nil {
			return nil, err
		}
		status = append(status, resp.Status)
	}
	return status, nil
}

func RegisterServices(node api.Node, context Context) (status []string, err error) {
	conf := context.Config

	for _, service := range getServices(conf["services"]) {
		resp, err := api.RegisterService(node, service, conf)
		if err != nil {
			return nil, err
		}
		status = append(status, resp.Status)
	}
	return status, nil
}
