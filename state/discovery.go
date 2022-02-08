package state

import (
	"fmt"
	"github.com/sav4enk0r0man/stolon-consul-discovery/api"
	"github.com/sav4enk0r0man/stolon-consul-discovery/config"
	"github.com/sav4enk0r0man/stolon-consul-discovery/logger"
	"strconv"
	"time"
)

type Context struct {
	Message string
	Error   error
	Config  config.Options
	Logger  logger.Logger
}

func (c Context) GetConf(key string) string {
	return c.Config[key]
}

func stringsContains(strs []string, substr string) bool {
	for _, s := range strs {
		if s == substr {
			return true
		}
	}
	return false
}

func Discovery(context chan Context) {
	var index int64 = 0
	var clusterState api.ClusterState

	ctx := <-context
	conf := ctx.Config
	log := ctx.Logger

	pollInterval, err := strconv.Atoi(conf["pollinterval"])
	if err != nil {
		context <- Context{
			Error:   err,
			Message: fmt.Sprintln("can't parse pollinterval"),
			Config:  conf,
		}
		return
	}

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)

		index, err = api.WaintIndex(index, ctx)
		if err != nil {
			log.Error.Println(err)
			continue
		}
		log.Info.Printf("current index: %d", index)

		newClusterState, err := getClusterState(ctx)
		if err != nil {
			log.Error.Println(err)
			continue
		}

		clusterData, err := api.ShowClusterData(ctx)
		if err != nil {
			log.Trace.Println(err)
		} else {
			log.Trace.Println(clusterData)
		}

		if newClusterState.Serialized() != clusterState.Serialized() {
			log.Info.Println("State changed")
			if len(GetMasters(newClusterState)) != 1 {
				log.Info.Println("inconsistent master node state in DCS, skip...")
				continue
			}
			if stringsContains(GetAddresses(newClusterState), "") {
				log.Info.Println("node without assigned address, skip...")
				continue
			}

			log.Info.Printf("Old cluster state: %s", clusterState.Serialized())
			log.Info.Printf("New cluster state: %s", newClusterState.Serialized())
			if len(clusterState.Nodes()) > 0 {
				var nodes []string
				for _, node := range append(newClusterState.Nodes(), clusterState.Nodes()...) {
					if !stringsContains(nodes, node.Name()) {
						nodes = append(nodes, node.Name())
					}
				}
				for _, n := range nodes {
					oldNodeState := clusterState.NodeByName(n)
					newNodeState := newClusterState.NodeByName(n)
					// log.Info.Printf("old cluster node: %v", oldNodeState.Serialized())
					// log.Info.Printf("new cluster node: %v", newNodeState.Serialized())
					if oldNodeState.Serialized() != newNodeState.Serialized() {
						if oldNodeState.Serialized() != "" {
							log.Info.Printf("Deregister services: %s for node: %s status: %v", conf["service"], n, Deregister(oldNodeState, conf["url"]).Status)
						}
						if newNodeState.Healthy() {
							log.Info.Printf("Registered services: %s for node: %s status: %v", conf["service"], n, Register(newNodeState, conf["service"], conf["url"]).Status)
						}
					}
				}
			} else if len(newClusterState.Nodes()) > 0 {
				for _, n := range newClusterState.Nodes() {
					log.Info.Printf("Deregister services: %s for node: %s status: %v", conf["service"], n.Name(), Deregister(n, conf["url"]).Status)
					log.Info.Printf("Registered services: %s for node: %s status: %v", conf["service"], n.Name(), Register(n, conf["service"], conf["url"]).Status)
				}
			}
			clusterState = newClusterState
		} else {
			log.Info.Printf("stolon cluster has not changed, current state: %s", clusterState.Serialized())
		}
		index += 1
	}
}
