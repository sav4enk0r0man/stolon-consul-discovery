package api

type Cluster struct {
	Uid        string `json:"uid"`
	Generation int    `json:"generation"`
	ChangeTime string `json:"changeTime"`
	Spec       struct {
		SynchronousReplication           bool   `json:"synchronousReplication"`
		AdditionalWalSenders             string `json:"additionalWalSenders"`
		AdditionalMasterReplicationSlots string `json:"additionalMasterReplicationSlots"`
		InitMode                         string `json:"initMode"`
		NewConfig                        struct {
			DataChecksums bool `json:"dataChecksums"`
		} `json:"newConfig"`
		PgParameters PgParameters `json:"pgParameters"`
		// "pgHBA": null,
		// "automaticPgRestart": null
	} `json:"spec"`
	Status struct {
		Phase  string `json:"phase"`
		Master string `json:"master"`
	} `json:"status"`
}

type Keeper struct {
	Uid        string   `json:"uid"`
	Generation int      `json:"generation"`
	ChangeTime string   `json:"changeTime"`
	Spec       struct{} `json:"spec"`
	Status     struct {
		Healthy               bool   `json:"healthy"`
		LastHealthyTime       string `json:"lastHealthyTime"`
		BootUUID              string `json:"bootUUID"`
		PostgresBinaryVersion struct {
			Maj int `json:"Maj"`
			Min int `json:"Min"`
		} `json:"postgresBinaryVersion"`
	} `json:"status"`
}

type DB struct {
	Uid        string `json:"uid"`
	Generation int    `json:"generation"`
	ChangeTime string `json:"changeTime"`
	Spec       struct {
		KeeperUID            string `json:"keeperUID"`
		RequestTimeout       string `json:"requestTimeout"`
		MaxStandbys          int    `json:"maxStandbys"`
		AdditionalWalSenders int    `json:"additionalWalSenders"`
		// "additionalReplicationSlots": null,
		InitMode     string       `json:"initMode"`
		PgParameters PgParameters `json:"pgParameters"`
		// "pgHBA": null,
		Role         string `json:"role,omitempty"`
		FollowConfig struct {
			Type  string `json:"type"`
			Dbuid string `json:"dbuid"`
		} `json:"followConfig,omitempty"`
		Followers                   []string `json:"followers,omitempty"`
		SynchronousStandbys         []string `json:"synchronousStandbys,omitempty"`
		ExternalSynchronousStandbys []string `json:"externalSynchronousStandbys,omitempty"`
	} `json:"spec"`
	Status struct {
		Healthy           bool   `json:"healthy"`
		CurrentGeneration int    `json:"currentGeneration"`
		ListenAddress     string `json:"listenAddress"`
		Port              string `json:"port"`
		SystemdID         string `json:"systemdID"`
		TimelineID        int    `json:"timelineID"`
		XLogPos           int64  `json:"xLogPos"`

		PgParameters PgParameters `json:"pgParameters"`
		// "pgHBA": null,
		Role                        string   `json:"role,omitempty"`
		Followers                   []string `json:"followers,omitempty"`
		SynchronousStandbys         []string `json:"synchronousStandbys,omitempty"`
		ExternalSynchronousStandbys []string `json:"externalSynchronousStandbys,omitempty"`
		OlderWalFile                string   `json:"olderWalFile,omitempty"`
	} `json:"status"`
}

type Proxy struct {
	Generation int    `json:"generation"`
	ChangeTime string `json:"changeTime"`
	Spec       struct {
		MasterDbUid    string   `json:"masterDbUid"`
		EnabledProxies []string `json:"enabledProxies"`
	} `json:"spec"`
	Status struct{} `json:"status"`
}

type ClusterData struct {
	FormatVersion int               `json:"formatVersion"`
	ChangeTime    string            `json:"changeTime"`
	Cluster       Cluster           `json:"cluster"`
	Keepers       map[string]Keeper `json:"keepers,omitempty"`
	Dbs           map[string]DB     `json:"dbs,omitempty"`
	Proxy         Proxy             `json:"proxy,omitempty"`
}

type PgParameters struct {
	Datestyle               string `json:"datestyle,omitempty"`
	DefaultTextSearchConfig string `json:"default_text_search_config,omitempty"`
	DynamicSharedMemoryType string `json:"dynamic_shared_memory_type,omitempty"`
	LcMessages              string `json:"lc_messages,omitempty"`
	LcMonetary              string `json:"lc_monetary,omitempty"`
	LcNumeric               string `json:"lc_numeric,omitempty"`
	LcTime                  string `json:"lc_time,omitempty"`
	LogDestination          string `json:"log_destination,omitempty"`
	LogDirectory            string `json:"log_directory,omitempty"`
	LogFilename             string `json:"log_filename,omitempty"`
	LogLinePrefix           string `json:"log_line_prefix,omitempty"`
	LogRotationAge          string `json:"log_rotation_age,omitempty"`
	LogRotationSize         string `json:"log_rotation_size,omitempty"`
	LogTimezone             string `json:"log_timezone,omitempty"`
	LogTruncateOnRotation   string `json:"log_truncate_on_rotation,omitempty"`
	LoggingCollector        string `jsonf:"logging_collector,omitempty"`
	MaxConnections          string `json:"max_connections,omitempty"`
	MaxWalSize              string `json:"max_wal_size,omitempty"`
	MinWalSize              string `json:"min_wal_size,omitempty"`
	SharedBuffers           string `json:"shared_buffers,omitempty"`
	Timezone                string `json:"timezone,omitempty"`
	WalLevel                string `json:"wal_level,omitempty"`
}
