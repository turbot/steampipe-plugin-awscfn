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

func tableCloudformationOutput(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudformation_output",
		Description: "Cloudformation resource information",
		List: &plugin.ListConfig{
			Hydrate:    listCloudformationOutputs,
			KeyColumns: plugin.OptionalColumns([]string{"path"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "Resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "Specifies the resource properties.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "description",
				Description: "Output description.",
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

type cloudformationOutput struct {
	Name        string
	Value       interface{}
	Description interface{}
	StartLine   int
	Path        string
}

type OutputStruct struct {
	Outputs map[string]interface{} `cty:"Outputs"`
}

func listCloudformationOutputs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
			plugin.Logger(ctx).Error("cloudformation_output.listCloudformationOutputs", "file_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file %s: %v", path, err)
		}

		// Parse file contents
		var body interface{}
		if err := yaml.Unmarshal(content, &body); err != nil {
			panic(err)
		}
		body = convert(body)

		var result OutputStruct
		if b, err := json.Marshal(body); err != nil {
			panic(err)
		} else {
			err = json.Unmarshal(b, &result)
			if err != nil {
				plugin.Logger(ctx).Error("cloudformation_output.listCloudformationOutputs", "parse_error", err, "path", path)
				return nil, fmt.Errorf("failed to unmarshal file content %s: %v", path, err)
			}
		}

		// Decode file contents
		var root yaml.Node
		r := bytes.NewReader(content)
		decoder := yaml.NewDecoder(r)
		err = decoder.Decode(&root)
		if err != nil {
			plugin.Logger(ctx).Error("cloudformation_output.listCloudformationOutputs", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse file: %v", err)
		}
		var rows Rows
		treeToList(&root, []string{}, &rows, "Outputs")

		for k, v := range result.Outputs {
			test := v.(map[string]interface{})
			var lineNo int
			for _, r := range rows {
				if r.Name == k {
					lineNo = r.StartLine
				}
			}

			d.StreamListItem(ctx, cloudformationOutput{
				Name:        k,
				Value:       test["Value"],
				Description: test["Description"],
				StartLine:   lineNo,
				Path:        path,
			})
		}
	}

	return nil, nil
}
