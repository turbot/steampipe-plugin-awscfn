package main

import (
	"github.com/turbot/steampipe-plugin-awscfn/awscfn"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: awscfn.Plugin})
}
