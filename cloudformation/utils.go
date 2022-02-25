package cloudformation

import (
	"strconv"

	"gopkg.in/yaml.v3"
)

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

func treeToList(tree *yaml.Node, prefix []string, rows *Rows) {
	switch tree.Kind {
	case yaml.DocumentNode:
		for _, v := range tree.Content {
			treeToList(v, prefix, rows)
		}
	case yaml.SequenceNode:
		for i, v := range tree.Content {
			newKey := make([]string, len(prefix))
			copy(newKey, prefix)
			newKey = append(newKey, strconv.Itoa(i))
			treeToList(v, newKey, rows)
		}
	case yaml.MappingNode:
		if len(prefix) == 2 && prefix[0] == "Resources" {
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
			treeToList(val, newKey, rows)
		}
	case yaml.ScalarNode:
	}
}
