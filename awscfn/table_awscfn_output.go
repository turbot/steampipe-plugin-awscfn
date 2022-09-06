package awscfn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"gopkg.in/yaml.v3"
)

func tableAWSCFNOutput(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "awscfn_output",
		Description: "CloudFormation resource information",
		List: &plugin.ListConfig{
			Hydrate:    listAWSCloudFormationOutputs,
			KeyColumns: plugin.OptionalColumns([]string{"path"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "An identifier for the current output.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "The value of the property returned by the aws cloudformation describe-stacks command. The value of an output can include literals, parameter references, pseudo-parameters, a mapping value, or intrinsic functions.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "A String type that describes the output value. The value for the description declaration must be a literal string that's between 0 and 1024 bytes in length. You can't use a parameter or function to specify the description. The description can be a maximum of 4 K in length.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "export",
				Description: "The name of the resource output to be exported for a cross-stack reference.",
				Type:        proto.ColumnType_JSON,
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

type awsCFNOutput struct {
	Name        string
	Value       interface{}
	Description interface{}
	Export      interface{}
	StartLine   int
	Path        string
}

type OutputStruct struct {
	Outputs   map[string]interface{} `cty:"Outputs"`
	Resources map[string]interface{} `cty:"Resources"`
}

func listAWSCloudFormationOutputs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
			plugin.Logger(ctx).Error("awscfn_output.listAWSCloudFormationOutputs", "file_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Parse file contents
		var body interface{}
		if err := yaml.Unmarshal(content, &IncludeProcessor{&body}); err != nil {
			panic(err)
		}
		body = convert(body)

		var result OutputStruct
		if b, err := json.Marshal(body); err != nil {
			panic(err)
		} else {
			err = json.Unmarshal(b, &result)
			if err != nil {
				plugin.Logger(ctx).Error("awscfn_output.listAWSCloudFormationOutputs", "parse_error", err, "path", path)
				return nil, fmt.Errorf("failed to unmarshal file content %s: %w", path, err)
			}
		}

		// Fail if no Resources attribute defined in template file
		if result.Resources == nil {
			plugin.Logger(ctx).Error("awscfn_output.listAWSCloudFormationOutputs", "template_format_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse AWS CloudFormation template from file %s: Template format error: At least one Resources member must be defined", path)
		}

		// Decode file contents
		var root yaml.Node
		r := bytes.NewReader(content)
		decoder := yaml.NewDecoder(r)
		err = decoder.Decode(&root)
		if err != nil {
			plugin.Logger(ctx).Error("awscfn_output.listAWSCloudFormationOutputs", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse file: %w", err)
		}
		var rows Rows
		treeToList(&root, []string{}, &rows, "Outputs")

		for k, v := range result.Outputs {
			data := v.(map[string]interface{})

			// Return error, if Outputs map has missing Value defined
			if data["Value"] == nil {
				plugin.Logger(ctx).Error("awscfn_output.listAWSCloudFormationOutputs", "template_format_error", err, "path", path)
				return nil, fmt.Errorf("failed to parse AWS CloudFormation template from file %s: Template format error: Every Outputs member must contain a Value object with non-null value", path)
			}

			var lineNo int
			for _, r := range rows {
				if r.Name == k {
					lineNo = r.StartLine
				}
			}

			d.StreamListItem(ctx, awsCFNOutput{
				Name:        k,
				Value:       data["Value"],
				Description: data["Description"],
				Export:      data["Export"],
				StartLine:   lineNo,
				Path:        path,
			})
		}
	}

	return nil, nil
}
