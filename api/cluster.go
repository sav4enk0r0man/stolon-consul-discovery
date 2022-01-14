package api

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/client"
)

func ClusterState(clusterName string, url string) string {
	api := fmt.Sprintf("%s/v1/kv/stolon/cluster/%s/clusterdata", url, clusterName)
	response := client.ConsulClient(api)

	indexes := make([]Index, 0)
	err := json.Unmarshal(response, &indexes)
	if err != nil {
		fmt.Printf("unmarshal error: %s", err)
	}

	decoded, err := b64.URLEncoding.DecodeString(string(indexes[0].Value))
	if err != nil {
		fmt.Printf("decode error: %s", err)
	}
	return string(decoded)
}
