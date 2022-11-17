package awscfn

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type awscfnConfig struct {
	Paths []string `cty:"paths" steampipe:"watch"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"paths": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
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
