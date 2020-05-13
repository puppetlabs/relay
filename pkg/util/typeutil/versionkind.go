package typeutil

import "gopkg.in/yaml.v2"

type VersionKind struct {
	APIVersion string `yaml:"apiVersion" json:"apiVersion,omitempty"`
	Kind       string `yaml:"kind" json:"kind,omitempty"`
}

type VersionKindExpectation struct {
	APIVersion string
	Kind       string
}

func (vke *VersionKindExpectation) NewFromYAMLString(data string) (*VersionKind, error) {
	vk := &VersionKind{}
	if err := yaml.Unmarshal([]byte(data), &vk); err != nil {
		if _, ok := err.(*yaml.TypeError); ok {
			return nil, &InvalidVersionKindError{
				ExpectedVersion: vke.APIVersion,
				ExpectedKind:    vke.Kind,
			}
		}
		return nil, err
	}

	if vk.APIVersion != vke.APIVersion || vk.Kind != vke.Kind {
		return nil, &InvalidVersionKindError{
			ExpectedVersion: vke.APIVersion,
			ExpectedKind:    vke.Kind,
			GotVersion:      vk.APIVersion,
			GotKind:         vk.Kind,
		}
	}

	return vk, nil
}

func NewVersionKindExpectation(apiVersion, kind string) *VersionKindExpectation {
	return &VersionKindExpectation{
		APIVersion: apiVersion,
		Kind:       kind,
	}
}
