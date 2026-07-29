package parsers

import "gopkg.in/yaml.v3"

type yamlParser struct{}

func (p yamlParser) parse(data []byte) (map[string]any, error) {
	var root any
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	return rootMap(root)
}
