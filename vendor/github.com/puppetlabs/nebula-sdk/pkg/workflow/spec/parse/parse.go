package parse

import (
	"os"
	"path/filepath"
	"strings"
)

type Tree interface{}

func ParseFile(name string) (Tree, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch strings.ToLower(filepath.Ext(f.Name())) {
	case "yaml", "yml":
		return ParseYAML(f)
	default:
		return ParseJSON(f)
	}
}
