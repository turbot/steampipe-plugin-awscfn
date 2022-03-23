package awscfn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/awslabs/goformation/v6"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
	"gopkg.in/yaml.v3"
)

func tableAWSCFNResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "awscfn_resource",
		Description: "Cloudformation resource information.",
		List: &plugin.ListConfig{
			Hydrate:    listAWSCloudFormationResources,
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
				Name:        "properties_src",
				Description: "Specifies the resource properties defined in the template.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LiteralValue"),
			},
			{
				Name:        "properties",
				Description: "Specifies the resource properties with calculated values as per given condition or parameter reference.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "condition",
				Description: "Specifies the resource conditions.",
				Type:        proto.ColumnType_STRING,
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

type awsCFNResource struct {
	Name                string
	StartLine           int
	Type                string
	Path                string
	LiteralValue        interface{}
	Properties          interface{}
	CreationPolicy      interface{}
	DeletionPolicy      interface{}
	DependsOn           interface{}
	Metadata            interface{}
	UpdatePolicy        interface{}
	UpdateReplacePolicy interface{}
}

type resourceStruct struct {
	Resources map[string]interface{} `cty:"Resources"`
}

type templateStruct struct {
	Properties interface{} `json:"Properties"`
}

func listAWSCloudFormationResources(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
			plugin.Logger(ctx).Error("awscfn_resource.listAWSCloudFormationResources", "file_error", err, "path", path)
			return nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Parse file contents
		var body interface{}
		content = formatFileContent(content)
		if err := yaml.Unmarshal(content, &IncludeProcessor{&body}); err != nil {
			panic(err)
		}
		body = convert(body)

		var resourceData resourceStruct
		if b, err := json.Marshal(body); err != nil {
			panic(err)
		} else {
			err = json.Unmarshal(b, &resourceData)
			if err != nil {
				plugin.Logger(ctx).Error("awscfn_resource.listAWSCloudFormationResources", "parse_error", err, "path", path)
				return nil, fmt.Errorf("failed to unmarshal file content %s: %w", path, err)
			}
		}

		// Fail if no Resources attribute defined in template file
		if resourceData.Resources == nil {
			plugin.Logger(ctx).Error("awscfn_resource.listAWSCloudFormationResources", "template_format_error", err, "path", path)
			return nil, fmt.Errorf("failed to parse AWS CloudFormation template from file %s: Template format error: At least one Resources member must be defined", path)
		}

		template, err := goformation.Open(path)
		if err != nil {
			plugin.Logger(ctx).Error("awscfn_resource.listAWSCloudFormationResources", "file_error", err, "path", path)
		}

		// Decode file contents
		var root yaml.Node
		r := bytes.NewReader(content)
		decoder := yaml.NewDecoder(r)
		err = decoder.Decode(&root)
		if err != nil {
			plugin.Logger(ctx).Error("awscfn_resource.listAWSCloudFormationResources", "parse_error", err, "path", path)
			return nil, fmt.Errorf("failed to decode file content: %w", err)
		}
		var rows Rows
		treeToList(&root, []string{}, &rows, "Resources")

		for k, v := range resourceData.Resources {
			data := v.(map[string]interface{})

			// Return error, if Resources map has missing Type defined
			if data["Type"] == nil {
				plugin.Logger(ctx).Error("awscfn_resource.listAWSCloudFormationResources", "template_format_error", err, "path", path)
				return nil, fmt.Errorf("failed to parse AWS CloudFormation template from file %s: Template format error: Every Resources object must contain a Type member. Resource: %s", path, k)
			}

			// Return error if Properties defined with no value, or null
			_, isPresent := data["Properties"]
			if isPresent && data["Properties"] == nil {
				plugin.Logger(ctx).Error("awscfn_resource.listAWSCloudFormationResources", "template_format_error", err, "path", path)
				return nil, fmt.Errorf("[/Resources/%s/Properties] 'null' values are not allowed in templates. File: %s", k, path)
			}

			var lineNo int
			for _, r := range rows {
				if r.Name == k {
					lineNo = r.StartLine
				}
			}

			var propertyValue interface{}
			if template != nil {
				for mapKey, val := range template.Resources {
					if k == mapKey {
						reqBodyBytes := new(bytes.Buffer)
						err := json.NewEncoder(reqBodyBytes).Encode(val)
						if err != nil {
							plugin.Logger(ctx).Error("awscfn_resource.listAWSCloudFormationResources", "parse_error", err, "path", path)
							return nil, fmt.Errorf("failed to encode file content %s: %w", path, err)
						}

						byteData := reqBodyBytes.String()
						var result templateStruct
						err = json.Unmarshal([]byte(byteData), &result)
						if err != nil {
							plugin.Logger(ctx).Error("awscfn_resource.listAWSCloudFormationResources", "parse_error", err, "path", path)
							return nil, fmt.Errorf("failed to unmarshal resource content: %w", err)
						}
						propertyValue = result.Properties
					}
				}
			}

			d.StreamListItem(ctx, awsCFNResource{
				Name:                k,
				StartLine:           lineNo,
				Type:                data["Type"].(string),
				Path:                path,
				LiteralValue:        data["Properties"],
				Properties:          propertyValue,
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
