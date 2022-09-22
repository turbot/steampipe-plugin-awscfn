package awscfn

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"gopkg.in/yaml.v3"
)

// List all files as per configured paths
func listFilesByPath(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	awscfnConfig := GetConfig(d.Connection)
	paths := awscfnConfig.Paths
	if paths == nil {
		return nil, errors.New("paths must be configured")
	}

	var matches []string
	for _, i := range paths {
		// List the files in the given source directory
		files, err := d.GetSourceFiles(i)
		if err != nil {
			return nil, err
		}
		plugin.Logger(ctx).Warn("listFilesByPath", "path", i, "matches", files)
		matches = append(matches, files...)
	}

	plugin.Logger(ctx).Warn("listFilesByPath", "matches", matches)

	// Sanitize the matches to likely cloudformation files
	var fileList []string
	for _, i := range matches {
		// Check if file or directory
		fileInfo, err := os.Stat(i)
		if err != nil {
			plugin.Logger(ctx).Error("utils.listFilesByPath", "error getting file info", err, "path", i)
			return nil, err
		}

		// Ignore directories
		if fileInfo.IsDir() {
			continue
		}

		hit := false
		for _, j := range paths {
			if i == j {
				hit = true
				break
			}
		}
		if hit {
			fileList = append(fileList, i)
			continue
		}
		fileList = append(fileList, i)
	}
	return fileList, nil
}

func convert(i interface{}) interface{} {
	switch valueType := i.(type) {
	case map[interface{}]interface{}:
		data := map[string]interface{}{}
		for k, v := range valueType {
			data[k.(string)] = convert(v)
		}
		return data
	case []interface{}:
		for i, v := range valueType {
			valueType[i] = convert(v)
		}
	}
	return i
}

type Rows []Row
type Row struct {
	Name      string
	StartLine int
}

func treeToList(tree *yaml.Node, prefix []string, rows *Rows, searchObjectName string) {
	switch tree.Kind {
	case yaml.DocumentNode:
		for _, v := range tree.Content {
			treeToList(v, prefix, rows, searchObjectName)
		}
	case yaml.SequenceNode:
		if len(tree.Content) > 0 {
			row := Row{
				Name:      strings.Join(prefix, "."),
				StartLine: tree.Line,
			}
			*rows = append(*rows, row)
		}
	case yaml.MappingNode:
		if len(prefix) == 2 && prefix[0] == searchObjectName {
			row := Row{
				Name:      prefix[1],
				StartLine: tree.Line,
			}
			*rows = append(*rows, row)
		}
		i := 0
		for i < len(tree.Content)-1 {
			key := tree.Content[i]
			val := tree.Content[i+1]
			i = i + 2
			newKey := make([]string, len(prefix))
			copy(newKey, prefix)
			newKey = append(newKey, key.Value)
			treeToList(val, newKey, rows, searchObjectName)
		}
	case yaml.ScalarNode:
		if len(prefix) > 0 && prefix[0] == searchObjectName {
			row := Row{
				Name:      strings.Join(prefix, "."),
				StartLine: tree.Line,
			}
			*rows = append(*rows, row)
		}
	}
}

type Fragment struct {
	content *yaml.Node
}

func (f *Fragment) UnmarshalYAML(value *yaml.Node) error {
	var err error
	f.content, err = resolveCustomTags(value)
	return err
}

type IncludeProcessor struct {
	target interface{}
}

func (i *IncludeProcessor) UnmarshalYAML(value *yaml.Node) error {
	resolved, err := resolveCustomTags(value)
	if err != nil {
		return err
	}
	return resolved.Decode(i.target)
}

// resolveCustomTags preserves YAML short tags while parsing
func resolveCustomTags(node *yaml.Node) (*yaml.Node, error) {
	switch node.Tag {
	case "!Base64":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::Base64: %v", node.Value)), &f)
		return f.content, err
	case "!Cidr":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::Cidr: %v", node.Value)), &f)
		return f.content, err
	case "!GetAtt":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		if strings.Contains(node.Value, ".") {
			node.Value = fmt.Sprintf("%s", strings.Split(node.Value, "."))
		}
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::GetAtt: %v", node.Value)), &f)
		return f.content, err
	case "!GetAZs":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::GetAZs: %v", node.Value)), &f)
		return f.content, err
	case "!ImportValue":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::ImportValue:: %v", node.Value)), &f)
		return f.content, err
	case "!Join":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::Join: %v", node.Value)), &f)
		return f.content, err
	case "!Select":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::Select: %v", node.Value)), &f)
		return f.content, err
	case "!Split":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::Split: %v", node.Value)), &f)
		return f.content, err
	case "!Sub":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::Sub: %v", node.Value)), &f)
		return f.content, err
	case "!Transform":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::Transform: %v", node.Value)), &f)
		return f.content, err
	case "!Ref":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Ref: %v", node.Value)), &f)
		return f.content, err
	case "!Condition":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Condition: %v", node.Value)), &f)
		return f.content, err
	}
	if node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode {
		var err error
		for i := range node.Content {
			node.Content[i], err = resolveCustomTags(node.Content[i])
			if err != nil {
				return nil, err
			}
		}
	}
	return node, nil
}

func formatFileContent(content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte("!If"), []byte(fmt.Sprintf("\n%sFn::If:", strings.Repeat(" ", 8))))
	content = bytes.ReplaceAll(content, []byte("!Equals"), []byte(fmt.Sprintf("\n%sFn::Equals:", strings.Repeat(" ", 8))))
	content = bytes.ReplaceAll(content, []byte("!FindInMap"), []byte(fmt.Sprintf("\n%sFn::FindInMap:", strings.Repeat(" ", 8))))
	return content
}
