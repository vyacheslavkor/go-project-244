package parsers

import "encoding/json"

type jsonParser struct{}

func (p jsonParser) parse(data []byte) (map[string]any, error) {
	var root any
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	return rootMap(root)
}
