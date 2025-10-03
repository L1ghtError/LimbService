include .env
export

.PHONY: openapi
openapi: openapi_http

.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh infra ./internal/ports ports

.PHONY: lint
lint:
	@./scripts/lint.sh

.PHONY: fmt
fmt:
	goimports -l -w internal/