package parsers

import "gopkg.in/yaml.v3"

type yamlParser struct{}

func (p yamlParser) parse(data []byte) (parsed map[string]any, err error) {
	err = yaml.Unmarshal(data, &parsed)
	return parsed, err
}
