package awscfn

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-awscfn"

// Plugin creates this (awscfn) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"awscfn_mapping":   tableAWSCFNMapping(ctx),
			"awscfn_output":    tableAWSCFNOutput(ctx),
			"awscfn_parameter": tableAWSCFNParameter(ctx),
			"awscfn_resource":  tableAWSCFNResource(ctx),
		},
	}

	return p
}
