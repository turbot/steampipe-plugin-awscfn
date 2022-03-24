package awscfn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"gopkg.in/yaml.v3"
)

func tableAWSCFNMapping(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "awscfn_mapping",
		Description: "CloudFormation mapping information.",
		List: &plugin.ListConfig{
			Hydrate:    listAWSCloudFormationMappings,
			KeyColumns: plugin.OptionalColumns([]string{"path"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "map",
				Description: "Mapping name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key",
				Description: "The key name that maps to name-value pairs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "The name from the name-value pair.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "The value from the name-value pair.",
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

type awsCFNMapping struct {
	Map       string
	Key       string
	Name      string
	Value     string
	StartLine int
	Path      string
}

type MappingsStruct struct {
	Mappings  map[string]interface{} `cty:"Mappings"`
	Resources map[string]interface{} `cty:"Resources"`
}

func listAWSCloudFormationMappings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
			plugin.Logger(ctx).Error("awscfn_mapping.listAWSCloudFormationMappings", "file_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Parse file contents
		var body interface{}
		if err := yaml.Unmarshal(content, &body); err != nil {
			panic(err)
		}
		body = convert(body)

		var result MappingsStruct
		if b, err := json.Marshal(body); err != nil {
			panic(err)
		} else {
			err = json.Unmarshal(b, &result)
			if err != nil {
				plugin.Logger(ctx).Error("awscfn_mapping.listAWSCloudFormationMappings", "parse_error", err, "path", path)
				return nil, fmt.Errorf("failed to unmarshal file content %s: %w", path, err)
			}
		}

		// Fail if no Resources attribute defined in template file
		if result.Resources == nil {
			plugin.Logger(ctx).Error("awscfn_mapping.listAWSCloudFormationMappings", "template_format_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse AWS CloudFormation template from file %s: Template format error: At least one Resources member must be defined", path)
		}

		// Decode file contents
		var root yaml.Node
		r := bytes.NewReader(content)
		decoder := yaml.NewDecoder(r)
		err = decoder.Decode(&root)
		if err != nil {
			plugin.Logger(ctx).Error("awscfn_mapping.listAWSCloudFormationMappings", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse file: %w", err)
		}
		var rows Rows
		treeToList(&root, []string{}, &rows, "Mappings")

		for k, v := range result.Mappings {
			data := v.(map[string]interface{})
			// TODO: Fix line numbers to represent the start of each name-value pair instead of the map
			var lineNo int
			for _, r := range rows {
				if r.Name == k {
					lineNo = r.StartLine
				}
			}

			for mapKey, mapValue := range data {
				for nameKey, nameValue := range mapValue.(map[string]interface{}) {
					d.StreamListItem(ctx, awsCFNMapping{
						Map:       k,
						Key:       mapKey,
						Name:      nameKey,
						Value:     nameValue.(string),
						StartLine: lineNo,
						Path:      path,
					})
				}
			}
		}
	}

	return nil, nil
}
