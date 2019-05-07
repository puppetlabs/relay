package workflow

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type Variable struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

type Action struct {
	Name       string     `yaml:"name" json:"name"`
	Image      string     `yaml:"image" json:"image"`
	ResourceID string     `yaml:"resourceID" json:"resource_id"`
	Spec       ActionSpec `yaml:"spec" json:"spec"`
}

type ActionSpec []byte

func (a *ActionSpec) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw map[string]interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	b, err := yaml.Marshal(raw)
	if err != nil {
		return err
	}

	*a = b

	return nil
}

type Workflow struct {
	APIVersion string      `yaml:"version" json:"api_version"`
	Variables  []*Variable `yaml:"variables" json:"variables"`
	Actions    []*Action   `yaml:"actions" json:"actions"`
	Steps      []string    `yaml:"steps" json:"steps"`
}

func (w *Workflow) Encode() ([]byte, error) {
	return json.Marshal(w)
}
