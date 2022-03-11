package awscloudformation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"gopkg.in/yaml.v3"
)

func tableAWSCloudFormationParameter(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "awscloudformation_parameter",
		Description: "Cloudformation parameter information",
		List: &plugin.ListConfig{
			Hydrate:    listAWSCloudFormationParameters,
			KeyColumns: plugin.OptionalColumns([]string{"path"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Parameter name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The data type for the parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_value",
				Description: "A value of the appropriate type for the template to use if no value is specified when a stack is created. If you define constraints for the parameter, you must specify a value that adheres to those constraints.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "max_length",
				Description: "An integer value that determines the largest number of characters you want to allow for String types.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "min_length",
				Description: "An integer value that determines the smallest number of characters you want to allow for String types.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "max_value",
				Description: "A numeric value that determines the largest numeric value you want to allow for Number types.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "min_value",
				Description: "A numeric value that determines the smallest numeric value you want to allow for Number types.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "no_echo",
				Description: "Indicates whether to mask the parameter value to prevent it from being displayed in the console, command line tools, or API. If you set the NoEcho attribute to true, CloudFormation returns the parameter value masked as asterisks (*****) for any calls that describe the stack or stack events, except for information stored in the locations specified below.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "allowed_pattern",
				Description: "A regular expression that represents the patterns to allow for String types. The pattern must match the entire parameter value provided.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allowed_values",
				Description: "An array containing the list of values allowed for the parameter.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "description",
				Description: "A string of up to 4000 characters that describes the parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "constraint_description",
				Description: "A string that explains a constraint when the constraint is violated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "start_line",
				Description: "Starting line number.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "path",
				Description: "Path to the file.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

type awsCloudFormationParameter struct {
	Name                  string
	Type                  string
	DefaultValue          interface{}
	Description           interface{}
	AllowedPattern        interface{}
	AllowedValues         interface{}
	ConstraintDescription interface{}
	MaxLength             interface{}
	MinLength             interface{}
	MaxValue              interface{}
	MinValue              interface{}
	NoEcho                interface{}
	StartLine             int
	Path                  string
}

type parametersStruct struct {
	Parameters map[string]interface{} `cty:"Parameters"`
	Resources  map[string]interface{} `cty:"Resources"`
}

func listAWSCloudFormationParameters(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// #1 - Path via qual
	// If the path was requested through qualifier then match it exactly. Globs
	// are not supported in this context since the output value for the column
	// will never match the requested value.
	//
	// #2 - Path via glob paths in config
	var paths []string
	if d.KeyColumnQuals["path"] != nil {
		paths = []string{d.KeyColumnQuals["path"].GetStringValue()}
	} else {
		var err error
		paths, err = listFilesByPath(ctx, d.Connection)
		if err != nil {
			return nil, err
		}
	}

	for _, path := range paths {
		// Read files
		content, err := ioutil.ReadFile(path)
		if err != nil {
			plugin.Logger(ctx).Error("awscloudformation_parameter.listAWSCloudFormationParameters", "file_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file %s: %v", path, err)
		}

		// Parse file contents
		var body interface{}
		if err := yaml.Unmarshal(content, &body); err != nil {
			panic(err)
		}
		body = convert(body)

		var result parametersStruct
		if b, err := json.Marshal(body); err != nil {
			panic(err)
		} else {
			err = json.Unmarshal(b, &result)
			if err != nil {
				plugin.Logger(ctx).Error("awscloudformation_parameter.listAWSCloudFormationParameters", "parse_error", err, "path", path)
				return nil, fmt.Errorf("failed to unmarshal file content %s: %v", path, err)
			}
		}

		// Fail if no Resources attribute defined in template file
		if result.Resources == nil {
			plugin.Logger(ctx).Error("awscloudformation_parameter.listAWSCloudFormationParameters", "template_format_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse AWS CloudFormation template from file %s: Template format error: At least one Resources member must be defined", path)
		}

		// Decode file contents
		var root yaml.Node
		r := bytes.NewReader(content)
		decoder := yaml.NewDecoder(r)
		err = decoder.Decode(&root)
		if err != nil {
			plugin.Logger(ctx).Error("awscloudformation_parameter.listAWSCloudFormationParameters", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse file: %v", err)
		}
		var rows Rows
		treeToList(&root, []string{}, &rows, "Parameters")

		for k, v := range result.Parameters {
			data := v.(map[string]interface{})

			// Return error, if Parameters map has missing Type defined
			if data["Type"] == nil {
				plugin.Logger(ctx).Error("awscloudformation_parameter.listAWSCloudFormationParameters", "template_format_error", err, "path", path)
				return nil, fmt.Errorf("failed to parse AWS CloudFormation template from file %s: Template format error: Every Parameters object must contain a Type member with non-null value", path)
			}

			var lineNo int
			for _, r := range rows {
				if r.Name == k {
					lineNo = r.StartLine
				}
			}
			d.StreamListItem(ctx, awsCloudFormationParameter{
				Name:                  k,
				Type:                  data["Type"].(string),
				DefaultValue:          data["Default"],
				Description:           data["Description"],
				AllowedPattern:        data["AllowedPattern"],
				AllowedValues:         data["AllowedValues"],
				ConstraintDescription: data["ConstraintDescription"],
				MaxLength:             data["MaxLength"],
				MinLength:             data["MinLength"],
				MaxValue:              data["MaxValue"],
				MinValue:              data["MinValue"],
				NoEcho:                data["NoEcho"],
				StartLine:             lineNo,
				Path:                  path,
			})
		}
	}

	return nil, nil
}
