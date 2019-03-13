package helm

import "fmt"

const defaultNamespace = "tiller-world"

type Chart struct {
	Path           string
	ValuesPaths    []string
	ValueOverrides map[string]string
	Namespace      string
}

func (c Chart) params() []string {
	var params []string

	namespace := c.Namespace
	if namespace == "" {
		namespace = defaultNamespace
	}

	params = append(params, "--namespace", namespace)

	for _, valuePath := range c.ValuesPaths {
		params = append(params, "-f", valuePath)
	}

	for k, v := range c.ValueOverrides {
		params = append(params, "--set", fmt.Sprintf("%s=%s", k, v))
	}

	params = append(params, c.Path)

	return params
}
