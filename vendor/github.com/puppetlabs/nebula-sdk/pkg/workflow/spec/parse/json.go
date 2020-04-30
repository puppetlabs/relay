package parse

import (
	"encoding/json"
	"io"
	"strings"
)

func ParseJSON(r io.Reader) (Tree, error) {
	var tree Tree
	if err := json.NewDecoder(r).Decode(&tree); err != nil {
		return nil, err
	}

	return tree, nil
}

func ParseJSONString(data string) (Tree, error) {
	return ParseJSON(strings.NewReader(data))
}
