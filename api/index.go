package api

import (
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/client"
)

type Index struct {
	LockIndex   int64  `json:"LockIndex"`
	Key         string `json:"Key"`
	Flags       int64  `json:"Flags"`
	Value       string `json:"Value"`
	CreateIndex int64  `json:"CreateIndex"`
	ModifyIndex int64  `json:"ModifyIndex"`
}

type IndexesResponse struct {
	Collection []Index
}

func WaintIndex(clusterName string, url string, index int64) []byte {
	api := fmt.Sprintf("%s/v1/kv/stolon/cluster/%s/clusterdata?wait=0s&index=%s",
		url, clusterName, fmt.Sprintf("%d", index))
	return client.ConsulClient(api)
}
