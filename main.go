package main

import (
	"github.com/sav4enk0r0man/stolon-consul-discovery/config"
	"github.com/sav4enk0r0man/stolon-consul-discovery/logger"
	"github.com/sav4enk0r0man/stolon-consul-discovery/state"
)

func main() {
	conf := config.Get()
	log := logger.NewLogger(conf)
	if err := config.Validate(conf); err != nil {
		log.Fatal.Println(err)
	}
	cluster := make(chan state.Context)

	log.Info.Println("Starting...")
	log.Debug.Printf("Config: %v", conf)

	if conf.IsSet("cluster") {
		log.Info.Printf("Single cluster configuration: %v", conf["cluster"])

		go state.Discovery(cluster)
		cluster <- state.Context{Config: conf, Logger: *log}
		status := <-cluster

		log.Info.Printf("Worker exited: %s", status.Message)
		if status.Error != nil {
			log.Error.Printf("%v\n", status.Error)
		}
		log.Debug.Printf("Config: %v", status.Config)

	} else if conf.IsSet("configdir") {
		log.Info.Printf("Multiple clusters config directory: %v", conf["configdir"])

		confFiles, err := config.WalkDir(conf["configdir"], conf["configmask"])
		if err != nil {
			log.Error.Fatal(err)
		}
		if len(confFiles) > 0 {
			log.Info.Printf("Config files: %v", confFiles)
			for _, confFile := range confFiles {
				clusterConf := config.Parse(confFile)
				clusterLog := logger.NewLogger(clusterConf)
				if err := config.Validate(clusterConf); err != nil {
					clusterLog.Fatal.Fatal(err)
				}
				go state.Discovery(cluster)
				cluster <- state.Context{Config: clusterConf, Logger: *clusterLog}
			}
			for range confFiles {
				status := <-cluster
				log.Info.Printf("Worker exited: %s", status.Message)
				if status.Error != nil {
					log.Error.Printf("%v\n", status.Error)
				}
				log.Debug.Printf("Config: %v", status.Config)
				log.Info.Printf("Waiting next workers...")
			}
		} else {
			log.Info.Printf("Cluster configuration files (%s/%s) not found...",
				conf["configdir"], conf["configmask"])
		}
	} else {
		log.Fatal.Fatalf("Cluster name(s) required. Set the -cluster or -config options")
	}
}
