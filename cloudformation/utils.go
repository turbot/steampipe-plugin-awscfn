package cloudformation

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
	cloudformationConfig := GetConfig(p)
	if cloudformationConfig.Paths == nil {
		return nil, errors.New("paths must be configured to query JSON files")
	}

	var matches []string
	for _, i := range cloudformationConfig.Paths {
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
	}
}
