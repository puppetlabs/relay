package integration

import (
	"io"

	"github.com/go-yaml/yaml"
	"github.com/puppetlabs/relay/pkg/debug"
)

type IntegrationConfig struct {
	//
	// Fields available across all versions of this configuration.
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`

	//
	// Fields that will likely be part of the first version.
	Name        string `yaml:"name"`
	Channel     string `yaml:"channel"`
	Version     string `yaml:"version"`
	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`
	License     string `yaml:"license"`
	Owner       struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
		URL   string `yaml:"url"`
	} `yaml:"owner"`

	Homepage string `yaml:"homepage"`
	Source   string `yaml:"source"`
}

func ReadConfig(r io.Reader) (*IntegrationConfig, error) {
	var config IntegrationConfig

	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		debug.Logf("failed to read integration configuration: %v", err)
		return nil, err
	}

	return &config, nil
}
