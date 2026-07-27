package parsers

import "encoding/json"

type jsonParser struct{}

func (p jsonParser) parse(data []byte) (parsed map[string]any, err error) {
	err = json.Unmarshal(data, &parsed)
	return parsed, err
}
