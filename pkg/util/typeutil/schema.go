package typeutil

import (
	"encoding/json"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

func ValidateYAMLString(schema *gojsonschema.Schema, data string) error {
	var root interface{}
	if err := yaml.Unmarshal([]byte(data), &root); err != nil {
		return err
	}

	b, err := json.Marshal(root)
	if err != nil {
		return err
	}

	result, err := schema.Validate(gojsonschema.NewBytesLoader(b))
	if err != nil {
		return err
	}

	return ValidationErrorFromResult(result)
}

func LoadSchemaFromStrings(primary string, rest ...string) (*gojsonschema.Schema, error) {
	loader := gojsonschema.NewSchemaLoader()
	loader.Validate = true

	for _, secondary := range rest {
		if err := loader.AddSchemas(gojsonschema.NewStringLoader(secondary)); err != nil {
			return nil, err
		}
	}

	return loader.Compile(gojsonschema.NewStringLoader(primary))
}
