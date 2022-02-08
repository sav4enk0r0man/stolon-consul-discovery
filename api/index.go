package api

import (
	"encoding/json"
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/client"
	"github.com/sav4enk0r0man/stolon-consul-discovery/logger"
)

type context interface {
	GetConf(string) string
}

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

func WaintIndex(index int64, ctx context) (int64, error) {
	clusterName := ctx.GetConf("cluster")
	url := ctx.GetConf("url")

	api := fmt.Sprintf("%s/v1/kv/stolon/cluster/%s/clusterdata?wait=0s&index=%s",
		url, clusterName, fmt.Sprintf("%d", index))

	response, err := client.Get(api, map[string]string{
		"httptimeout": ctx.GetConf("httptimeout"),
	})
	if err != nil {
		if err.Error() == "404" {
			return 0, logger.Wrapper(err, fmt.Sprintf("Cluster %s not found", clusterName))
		}
		return 0, err
	}

	indexes := make([]Index, 0)
	if err := json.Unmarshal(response, &indexes); err != nil {
		return 0, logger.Wrapper(err, err.Error())
	}

	return indexes[0].ModifyIndex, nil
}
