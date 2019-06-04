package workflow

import (
	"encoding/json"
)

type Variable struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

type Step struct {
	Name      string                 `yaml:"name" json:"name"`
	Image     string                 `yaml:"image" json:"image"`
	Spec      map[string]interface{} `yaml:"spec" json:"spec"`
	DependsOn string                 `yaml:"depends_on" json:"depends_on"`
}

type Workflow struct {
	Version     string      `yaml:"version" json:"version"`
	Description string      `yaml:"description" json:"description"`
	Variables   []*Variable `yaml:"variables" json:"variables"`
	Steps       []*Step     `yaml:"steps" json:"steps"`
}

func (w *Workflow) Encode() ([]byte, error) {
	return json.Marshal(w)
}
