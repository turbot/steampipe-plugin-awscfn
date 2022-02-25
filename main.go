package main

import (
	"github.com/turbot/steampipe-plugin-cloudformation/cloudformation"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: cloudformation.Plugin})
}
