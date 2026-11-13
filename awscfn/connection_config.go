package awscfn

import (
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type awscfnConfig struct {
	Paths []string `hcl:"paths,optional" steampipe:"watch"`
}

func ConfigInstance() interface{} {
	return &awscfnConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) awscfnConfig {
	if connection == nil || connection.GetConfig() == nil {
		return awscfnConfig{}
	}
	config, _ := connection.GetConfig().(awscfnConfig)
	return config
}
