package cloudformation

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

func tableCloudformationMapping(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudformation_mapping",
		Description: "Cloudformation mapping information",
		List: &plugin.ListConfig{
			Hydrate:    listCloudformationMappings,
			KeyColumns: plugin.OptionalColumns([]string{"path"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Parameter name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key",
				Description: "The data type for the parameter.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "An array containing the list of values allowed for the parameter.",
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

type cloudformationMapping struct {
	Name      string
	Key       string
	Value     interface{}
	StartLine int
	Path      string
}

type MappingsStruct struct {
	Mappings  map[string]interface{} `cty:"Mappings"`
	Resources map[string]interface{} `cty:"Resources"`
}

func listCloudformationMappings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
			plugin.Logger(ctx).Error("cloudformation_mapping.listCloudformationMappings", "file_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file %s: %v", path, err)
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
				plugin.Logger(ctx).Error("cloudformation_mapping.listCloudformationMappings", "parse_error", err, "path", path)
				return nil, fmt.Errorf("failed to unmarshal file content %s: %v", path, err)
			}
		}

		// Fail if no Resources attribute defined in template file
		if result.Resources == nil {
			plugin.Logger(ctx).Error("cloudformation_mapping.listCloudformationMappings", "template_format_error", err, "path", path)
			return nil, fmt.Errorf("Template format error: At least one Resources member must be defined. File: %s", path)
		}

		// Decode file contents
		var root yaml.Node
		r := bytes.NewReader(content)
		decoder := yaml.NewDecoder(r)
		err = decoder.Decode(&root)
		if err != nil {
			plugin.Logger(ctx).Error("cloudformation_mapping.listCloudformationMappings", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse file: %v", err)
		}
		var rows Rows
		treeToList(&root, []string{}, &rows, "Mappings")

		for k, v := range result.Mappings {
			data := v.(map[string]interface{})
			var lineNo int
			for _, r := range rows {
				if r.Name == k {
					lineNo = r.StartLine
				}
			}

			for mapKey, mapValue := range data {
				d.StreamListItem(ctx, cloudformationMapping{
					Name:      k,
					Key:       mapKey,
					Value:     mapValue,
					StartLine: lineNo,
					Path:      path,
				})
			}
		}
	}

	return nil, nil
}
