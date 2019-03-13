package helm

type Chart struct {
	Path        string
	ValuesPaths []string
}

func (c Chart) params() []string {
	var params []string

	for _, valuePath := range c.ValuesPaths {
		params = append(params, "-f", valuePath)
	}

	params = append(params, c.Path)

	return params
}
