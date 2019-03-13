package helm

import "fmt"

type Chart struct {
	Path           string
	ValuesPaths    []string
	ValueOverrides map[string]string
}

func (c Chart) params() []string {
	var params []string

	for _, valuePath := range c.ValuesPaths {
		params = append(params, "-f", valuePath)
	}

	for k, v := range c.ValueOverrides {
		params = append(params, "--set", fmt.Sprintf("%s=%s", k, v))
	}

	params = append(params, c.Path)

	return params
}
