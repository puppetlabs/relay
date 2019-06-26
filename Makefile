build:
	go build -mod=vendor ./cmd/nebula

client:
	api-spec-converter -f openapi_3 -t swagger_2 -s yaml ../nebula-api/openapi/swagger.yaml > swagger.yaml
	rm -rf pkg/client/{api,models}
	swagger generate client -f swagger.yaml -c pkg/client/api -m pkg/client/api/models --skip-validation
	rm -rf swagger.yaml