package api

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/client"
	"github.com/sav4enk0r0man/stolon-consul-discovery/logger"
	"sort"
)

type ClusterState struct {
	nodes      []Node
	serialized string
}

func NewClusterState() *ClusterState {
	return &ClusterState{}
}

func (s *ClusterState) AddNode(n Node) {
	s.nodes = append(s.nodes, n)
}

func (s *ClusterState) Nodes() []Node {
	return s.nodes
}

func (s *ClusterState) NodeByName(name string) Node {
	for _, node := range s.nodes {
		if name == node.name {
			return node
		}
	}
	return Node{}
}

func (s *ClusterState) Serialize() {
	sort.Slice(s.nodes, func(i, j int) bool {
		return s.nodes[i].Name() < s.nodes[j].Name()
	})
	for i := range s.Nodes() {
		s.nodes[i].Serialize()
		s.serialized += s.nodes[i].Serialized()
	}
}

func (s ClusterState) Serialized() string {
	return s.serialized
}

func GetClusterData(ctx context) ([]Index, error) {
	clusterName := ctx.GetConf("cluster")
	url := ctx.GetConf("url")

	api := fmt.Sprintf("%s/v1/kv/stolon/cluster/%s/clusterdata", url, clusterName)
	response, _ := client.Get(api, map[string]string{
		"httptimeout": ctx.GetConf("httptimeout"),
	})

	indexes := make([]Index, 0)
	err := json.Unmarshal(response, &indexes)
	if err != nil {
		return nil, logger.Wrapper(err, err.Error())
	}
	return indexes, nil
}

func ShowClusterData(ctx context) (string, error) {
	indexes, err := GetClusterData(ctx)
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(indexes, "", "    ")
	if err != nil {
		return "", logger.Wrapper(err, err.Error())
	}
	resp := fmt.Sprintf("Cluster raw data: %s\n", string(data))

	decodedValue, err := b64.URLEncoding.DecodeString(string(indexes[0].Value))
	if err != nil {
		return "", logger.Wrapper(err, err.Error())
	}

	var value ClusterData
	if err = json.Unmarshal([]byte(decodedValue), &value); err != nil {
		return "", logger.Wrapper(err, err.Error())
	}

	data, err = json.MarshalIndent(value, "", "    ")
	if err != nil {
		return "", logger.Wrapper(err, err.Error())
	}

	resp += fmt.Sprintf("\nCluster data value: %s\n", string(data))
	return resp, nil
}
