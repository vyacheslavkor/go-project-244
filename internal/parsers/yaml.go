package parsers

import "gopkg.in/yaml.v3"

type yamlParser struct{}

func (p yamlParser) parse(data []byte) (any, error) {
	var parsed any
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}
