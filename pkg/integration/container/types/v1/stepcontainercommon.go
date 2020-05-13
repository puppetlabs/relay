package v1

import (
	"gopkg.in/yaml.v3"
)

type FileSource string

var (
	FileSourceSystem FileSource = ""
	FileSourceSDK    FileSource = "sdk"
)

type FileRef struct {
	From FileSource
	Name string
}

func (fr *FileRef) UnmarshalYAML(value *yaml.Node) error {
	var ps string
	if err := value.Decode(&ps); err == nil {
		// Successfully decoded to string, so this is the simple case of a file
		// path.
		fr.Name = ps
		return nil
	} else if _, ok := err.(*yaml.TypeError); !ok {
		// An error other than a type error should just be returned.
		return err
	}

	// Now we try to decode into a struct instead.
	var pm struct {
		From FileSource `yaml:"from"`
		Name string     `yaml:"name"`
	}
	if err := value.Decode(&pm); err != nil {
		return err
	}

	fr.From = pm.From
	fr.Name = pm.Name
	return nil
}

type StepContainerImage struct {
	Template  FileRef  `yaml:"template"`
	DependsOn []string `yaml:"dependsOn"`
}

type StepContainerSetting struct {
	Description string
	Value       interface{}
}

func (scs *StepContainerSetting) UnmarshalYAML(value *yaml.Node) error {
	var setting interface{}
	if err := value.Decode(&setting); err != nil {
		return err
	}

	// We need to disambiguate to see if we should set the setting itself or
	// just the value.
	if m, ok := setting.(map[string]interface{}); ok {
		switch len(m) {
		case 2:
			description, ok := m["description"].(string)
			if !ok {
				break
			}

			value, found := m["value"]
			if !found {
				break
			}

			scs.Description = description
			scs.Value = value
			return nil
		case 1:
			if description, ok := m["description"].(string); ok {
				scs.Description = description
				return nil
			}

			if value, found := m["value"]; found {
				scs.Value = value
				return nil
			}
		}
	}

	// Nothing matched, so we assume the user meant to just assign the value.
	scs.Value = setting
	return nil
}

type StepContainerCommon struct {
	Inherit    *FileRef                        `yaml:"inherit,omitempty"`
	SDKVersion string                          `yaml:"sdkVersion"`
	Images     map[string]StepContainerImage   `yaml:"images"`
	Settings   map[string]StepContainerSetting `yaml:"settings"`
}
