package main

import (
	"github.com/turbot/steampipe-plugin-awscloudformation/awscloudformation"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: awscloudformation.Plugin})
}
