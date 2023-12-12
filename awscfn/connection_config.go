package awscfn

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type awscfnConfig struct {
	Paths []string `hcl:"paths" steampipe:"watch"`
}

func ConfigInstance() interface{} {
	return &awscfnConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) awscfnConfig {
	if connection == nil || connection.Config == nil {
		return awscfnConfig{}
	}
	config, _ := connection.Config.(awscfnConfig)
	return config
}
