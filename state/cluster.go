package state

import (
	b64 "encoding/base64"
	"encoding/json"
	"github.com/sav4enk0r0man/stolon-consul-discovery/api"
	"github.com/sav4enk0r0man/stolon-consul-discovery/logger"
)

func getClusterState(ctx Context) (api.ClusterState, error) {
	indexes, err := api.GetClusterData(ctx)
	if err != nil {
		return api.ClusterState{}, err
	}

	decoded, err := b64.URLEncoding.DecodeString(string(indexes[0].Value))
	if err != nil {
		return api.ClusterState{}, logger.Wrapper(err, err.Error())
	}

	var value api.ClusterData
	err = json.Unmarshal([]byte(decoded), &value)
	if err != nil {
		return api.ClusterState{}, logger.Wrapper(err, err.Error())
	}

	state := api.NewClusterState()
	for _, db := range value.Dbs {
		state.AddNode(*api.NewNode(db.Uid, db.Spec.KeeperUID,
			db.Spec.Role, db.Status.ListenAddress, db.Status.Healthy))
	}

	state.Serialize()
	return *state, nil
}
