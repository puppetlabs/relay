package openapi

//go:generate 	docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli:v5.1.1 generate -i https://api.relay.sh/openapi/latest -g go --global-property apis,models,supportingFiles,modelDocs=false -o /local/pkg/client/openapi
