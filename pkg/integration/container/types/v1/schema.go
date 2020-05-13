package v1

import (
	"fmt"

	"github.com/puppetlabs/relay/pkg/integration/container/asset"
	"github.com/puppetlabs/relay/pkg/util/typeutil"
	"github.com/xeipuuv/gojsonschema"
)

var (
	StepContainerTemplateSchema *gojsonschema.Schema
	StepContainerSchema         *gojsonschema.Schema
)

func init() {
	stepContainerTemplateSchema, err := typeutil.LoadSchemaFromStrings(
		asset.MustAssetString("schemas/v1/StepContainerTemplate.json"),
		asset.MustAssetString("schemas/v1/StepContainer-common.json"),
	)
	if err != nil {
		panic(fmt.Errorf("failed to load schema for StepContainerTemplate: %+v", err))
	}

	stepContainerSchema, err := typeutil.LoadSchemaFromStrings(
		asset.MustAssetString("schemas/v1/StepContainer.json"),
		asset.MustAssetString("schemas/v1/StepContainer-common.json"),
	)
	if err != nil {
		panic(fmt.Errorf("failed to load schema for StepContainer: %+v", err))
	}

	StepContainerTemplateSchema = stepContainerTemplateSchema
	StepContainerSchema = stepContainerSchema
}
