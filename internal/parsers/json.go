package parsers

import "encoding/json"

type jsonParser struct{}

func (p jsonParser) parse(data []byte) (any, error) {
	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}
