package v1

import (
	"io"
	"io/ioutil"

	"github.com/puppetlabs/relay/pkg/util/typeutil"
	"gopkg.in/yaml.v3"
)

const (
	StepContainerTemplateKind = "StepContainerTemplate"
)

var StepContainerTemplateVersionKindExpectation = typeutil.NewVersionKindExpectation(Version, StepContainerTemplateKind)

// StepContainerTemplate represents an object with kind "StepContainerTemplate".
type StepContainerTemplate struct {
	*typeutil.VersionKind `yaml:",inline"`
	*StepContainerCommon  `yaml:",inline"`
}

func NewStepContainerTemplateFromString(data string) (*StepContainerTemplate, error) {
	if _, err := StepContainerTemplateVersionKindExpectation.NewFromYAMLString(data); err != nil {
		return nil, err
	}

	if err := typeutil.ValidateYAMLString(StepContainerTemplateSchema, data); err != nil {
		return nil, err
	}

	sct := &StepContainerTemplate{}
	if err := yaml.Unmarshal([]byte(data), &sct); err != nil {
		return nil, err
	}

	return sct, nil
}

func NewStepContainerTemplateFromReader(r io.Reader) (*StepContainerTemplate, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return NewStepContainerTemplateFromString(string(b))
}
