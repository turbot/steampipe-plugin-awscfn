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

func tableCloudformationResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudformation_resource",
		Description: "Cloudformation resource information",
		List: &plugin.ListConfig{
			Hydrate:    listCloudformationResources,
			KeyColumns: plugin.OptionalColumns([]string{"path"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "An identifier for the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type identifies the type of resource that you are declaring.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "properties",
				Description: "Specifies the resource properties.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "creation_policy",
				Description: "Specifies the associated creation_policy with a resource to prevent its status from reaching create complete until AWS CloudFormation receives a specified number of success signals or the timeout period is exceeded.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "deletion_policy",
				Description: "With the deletion_policy attribute you can preserve, and in some cases, backup a resource when its stack is deleted. You specify a deletion_policy attribute for each resource that you want to control.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "depends_on",
				Description: "With the depends_on attribute you can specify that the creation of a specific resource follows another. When you add a depends_on attribute to a resource, that resource is created only after the creation of the resource specified in the depends_on attribute.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "metadata",
				Description: "The metadata attribute enables you to associate structured data with a resource. By adding a metadata attribute to a resource, you can add data in JSON or YAML to the resource declaration.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "update_policy",
				Description: "Use the update_policy attribute to specify how AWS CloudFormation handles updates to specific resources.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "update_replace_policy",
				Description: "Use the update_replace_policy attribute to retain or, in some cases, backup the existing physical instance of a resource when it's replaced during a stack update operation.",
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

type cloudformationResource struct {
	Name                string
	StartLine           int
	Type                string
	Path                string
	Properties          interface{}
	CreationPolicy      interface{}
	DeletionPolicy      interface{}
	DependsOn           interface{}
	Metadata            interface{}
	UpdatePolicy        interface{}
	UpdateReplacePolicy interface{}
}

type templateStruct struct {
	Resources map[string]interface{} `cty:"Resources"`
}

func listCloudformationResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
			plugin.Logger(ctx).Error("cloudformation_resource.listCloudformationResources", "file_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file %s: %v", path, err)
		}

		// Parse file contents
		var body interface{}
		if err := yaml.Unmarshal(content, &body); err != nil {
			panic(err)
		}
		body = convert(body)

		var result templateStruct
		if b, err := json.Marshal(body); err != nil {
			panic(err)
		} else {
			err = json.Unmarshal(b, &result)
			if err != nil {
				plugin.Logger(ctx).Error("cloudformation_resource.listCloudformationResources", "parse_error", err, "path", path)
				return nil, fmt.Errorf("failed to unmarshal file content %s: %v", path, err)
			}
		}

		// Decode file contents
		var root yaml.Node
		r := bytes.NewReader(content)
		decoder := yaml.NewDecoder(r)
		err = decoder.Decode(&root)
		if err != nil {
			plugin.Logger(ctx).Error("cloudformation_resource.listCloudformationResources", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse file: %v", err)
		}
		var rows Rows
		treeToList(&root, []string{}, &rows, "Resources")

		for k, v := range result.Resources {
			data := v.(map[string]interface{})
			var lineNo int
			for _, r := range rows {
				if r.Name == k {
					lineNo = r.StartLine
				}
			}
			d.StreamListItem(ctx, cloudformationResource{
				Name:                k,
				StartLine:           lineNo,
				Type:                data["Type"].(string),
				Path:                path,
				Properties:          data["Properties"],
				CreationPolicy:      data["CreationPolicy"],
				DeletionPolicy:      data["DeletionPolicy"],
				DependsOn:           data["DependsOn"],
				Metadata:            data["Metadata"],
				UpdatePolicy:        data["UpdatePolicy"],
				UpdateReplacePolicy: data["UpdateReplacePolicy"],
			})
		}
	}

	return nil, nil
}
