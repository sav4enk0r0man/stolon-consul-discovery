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

	if conf["deregister"] != "" {
		var msg string
		state, err := getClusterState(ctx)
		if err != nil {
			msg = fmt.Sprintln("can't get cluster state")
		} else {
			for _, node := range state.Nodes() {
				msg += fmt.Sprintf("\tderegister services: %v for node: %s status: ",
					getServices(conf["services"]), node.Name())
				status, err := DeregisterServices(node, ctx)
				if err != nil {
					break
				} else {
					msg += fmt.Sprintf("%v", status)
				}

			}
		}
		context <- Context{
			Error:   err,
			Message: msg,
			Config:  conf,
		}
		return
	}

	pollInterval, _ := strconv.Atoi(conf["pollinterval"])

DISCOVERY:
	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)

		index, err := api.WaintIndex(index, ctx)
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
			if len(getMasters(newClusterState)) != 1 {
				log.Info.Println("inconsistent master node state in DCS, skip...")
				continue
			}
			if stringsContains(getAddresses(newClusterState), "") {
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
					if oldNodeState.Serialized() != newNodeState.Serialized() {
						if oldNodeState.Serialized() != "" {
							msg := fmt.Sprintf("Deregister services: %s for node: %s status:", getServices(conf["services"]), n)
							status, err := DeregisterServices(oldNodeState, ctx)
							if err != nil {
								log.Error.Printf("%s %v", msg, err)
								continue DISCOVERY
							}
							log.Info.Printf("%s %s", msg, status)
						}
						if newNodeState.Healthy() {
							msg := fmt.Sprintf("Registered services: %s for node: %s status:", getServices(conf["services"]), n)
							status, err := RegisterServices(newNodeState, ctx)
							if err != nil {
								log.Error.Printf("%s %v", msg, err)
								continue DISCOVERY
							}
							log.Info.Printf("%s %s", msg, status)
						}
					}
				}
			} else if len(newClusterState.Nodes()) > 0 {
				for _, n := range newClusterState.Nodes() {
					msg := fmt.Sprintf("Deregister services: %s for node: %s status:", getServices(conf["services"]), n.Name())
					status, err := DeregisterServices(n, ctx)
					if err != nil {
						log.Error.Printf("%s %v", msg, err)
					} else {
						log.Info.Printf("%s %s", msg, status)
					}
					msg = fmt.Sprintf("Registered services: %s for node: %s status:", getServices(conf["services"]), n.Name())
					status, err = RegisterServices(n, ctx)
					if err != nil {
						log.Error.Printf("%s %v", msg, err)
						continue DISCOVERY
					}
					log.Info.Printf("%s %s", msg, status)
				}
			}
			clusterState = newClusterState
		} else {
			log.Info.Printf("stolon cluster has not changed, current state: %s", clusterState.Serialized())
		}
		index += 1
	}
}
