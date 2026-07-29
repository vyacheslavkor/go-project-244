package parsers

import "encoding/json"

type jsonParser struct{}

func (p jsonParser) parse(data []byte) (map[string]any, error) {
	var parsed map[string]any
	err := json.Unmarshal(data, &parsed)
	return parsed, err
}
