package awscloudformation

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type awsCloudFormationConfig struct {
	Paths []string `cty:"paths"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"paths": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
}

func ConfigInstance() interface{} {
	return &awsCloudFormationConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) awsCloudFormationConfig {
	if connection == nil || connection.Config == nil {
		return awsCloudFormationConfig{}
	}
	config, _ := connection.Config.(awsCloudFormationConfig)
	return config
}
