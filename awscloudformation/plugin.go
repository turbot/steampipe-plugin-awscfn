package awscloudformation

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-awscloudformation"

// Plugin creates this (awscloudformation) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"awscloudformation_mapping":   tableAWSCloudFormationMapping(ctx),
			"awscloudformation_output":    tableAWSCloudFormationOutput(ctx),
			"awscloudformation_parameter": tableAWSCloudFormationParameter(ctx),
			"awscloudformation_resource":  tableAWSCloudFormationResource(ctx),
		},
	}

	return p
}
