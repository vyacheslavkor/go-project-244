package parsers

import "gopkg.in/yaml.v3"

type yamlParser struct{}

func (p yamlParser) parse(data []byte) (map[string]any, error) {
	var parsed map[string]any
	err := yaml.Unmarshal(data, &parsed)
	return parsed, err
}
