package awscloudformation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"gopkg.in/yaml.v3"
)

func listFilesByPath(ctx context.Context, p *plugin.Connection) ([]string, error) {
	awsCloudFormationConfig := GetConfig(p)
	if awsCloudFormationConfig.Paths == nil {
		return nil, errors.New("paths must be configured")
	}

	var matches []string
	for _, i := range awsCloudFormationConfig.Paths {
		// Check to resolve ~ to home dir
		if strings.HasPrefix(i, "~") {
			// File system context
			home, err := os.UserHomeDir()
			if err != nil {
				plugin.Logger(ctx).Error("utils.listFilesByPath", "os.UserHomeDir error. ~ will not be expanded in paths", err, "path", i)
				return nil, fmt.Errorf("os.UserHomeDir error. ~ will not be expanded in paths")
			}

			// Resolve ~ to home dir
			if home != "" {
				if i == "~" {
					i = home
				} else if strings.HasPrefix(i, "~/") {
					i = filepath.Join(home, i[2:])
				}
			}
		}

		// Get full path
		fullPath, err := filepath.Abs(i)
		if err != nil {
			plugin.Logger(ctx).Error("utils.listFilesByPath", "invlaid path", err, "path", i)
			return nil, fmt.Errorf("failed to fetch absolute path: %s", i)
		}

		// Expand globs
		iMatches, err := doublestar.Glob(fullPath)
		if err != nil {
			plugin.Logger(ctx).Error("utils.listFilesByPath", "path is not a valid glob", err, "path", i)
			return matches, fmt.Errorf("path is not a valid glob: %s", i)
		}
		matches = append(matches, iMatches...)
	}

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
		for i, v := range tree.Content {
			newKey := make([]string, len(prefix))
			copy(newKey, prefix)
			newKey = append(newKey, strconv.Itoa(i))
			treeToList(v, newKey, rows, searchObjectName)
		}
	case yaml.MappingNode:
		if len(prefix) == 2 && prefix[0] == searchObjectName {
			row := Row{
				Name:      prefix[1],
				StartLine: (tree.Line - 1),
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
	}
}

// used for loading included files
type Fragment struct {
	content *yaml.Node
}

func (f *Fragment) UnmarshalYAML(value *yaml.Node) error {
	var err error
	// process includes in fragments
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
	case "!FindInMap":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::FindInMap: %v", node.Value)), &f)
		return f.content, err
	case "!GetAtt":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
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
	case "!And":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::And: %v", node.Value)), &f)
		return f.content, err
	case "!Equals":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::Equals: %v", node.Value)), &f)
		return f.content, err
	case "!If":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::If: %v", node.Value)), &f)
		return f.content, err
	case "!Not":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::Not: %v", node.Value)), &f)
		return f.content, err
	case "!Or":
		if node.Kind != yaml.ScalarNode {
			break
		}
		var f Fragment
		err := yaml.Unmarshal([]byte(fmt.Sprintf("Fn::Or: %v", node.Value)), &f)
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
