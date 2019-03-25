package types

import "gopkg.in/yaml.v2"

type Plan struct {
	Version   string      `yaml:"version"`
	Name      string      `yaml:"name"`
	Variables []*Variable `yaml:"variables"`
	Actions   []*Action   `yaml:"actions"`
	Workflows []*Workflow `yaml:"workflows"`
}

type Variable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Action struct {
	Name       string     `yaml:"name"`
	Image      string     `yaml:"image"`
	ResourceID string     `yaml:"resourceID"`
	Spec       ActionSpec `yaml:"spec"`
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
	Name        string   `yaml:"name"`
	Default     bool     `yaml:"default"`
	ActionNames []string `yaml:"actions"`
	Trigger     Trigger  `yaml:"trigger"`
}

type Trigger struct {
	On []*TriggerConditional `yaml:"on"`
}

type TriggerConditional struct {
	Event      string `yaml:"event"`
	TagSpec    string `yaml:"tag"`
	BranchSpec string `yaml:"branch"`
}
