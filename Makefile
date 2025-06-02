.PHONY: migrate-ent
migrate-ent:
	@echo "Running Ent migrations..."
	@go run cmd/migrate/main.go
	@echo "✅ Ent migrations complete"

.PHONY: generate-ent
generate-ent:
	@echo "Generating Ent schema..."
	@go run scripts/main.go -cmd generate-ent
	@echo "✅ Ent schema generated"

.PHONY: generate-keys
generate-keys:
	@echo "Generating encrypted RSA key pair for production..."
	@./scripts/generate-keys.sh
	@echo "✅ Encrypted key pair generated"

.PHONY: generate-dev-keys
generate-dev-keys:
	@echo "Generating unencrypted RSA key pair for development..."
	@./scripts/generate-dev-keys.sh
	@echo "✅ Development key pair generated"

.PHONY: install-hooks
install-hooks:
	@echo "Installing Git hooks..."
	@git config core.hooksPath .githooks
	@chmod +x .githooks/pre-commit
	@echo "✅ Git hooks installed"

.PHONY: run-hooks
run-hooks:
	@echo "Running Git hooks..."
	@.githooks/pre-commit

.PHONY: lint-fix
lint-fix:
	@echo "🧼 Running gofmt to auto-format..."
	gofmt -s -w .
	@echo "🧹 Running golangci-lint with --fix (autofix where possible)..."
	golangci-lint run --fix || true
	@echo "✅ Lint fixes applied (where possible)."


.PHONY: swagger-clean
swagger-clean:
	rm -rf docs/swagger

.PHONY: install-swag
install-swag:
	@which swag > /dev/null || (go install github.com/swaggo/swag/cmd/swag@latest)

.PHONY: swagger
swagger: swagger-2-0 swagger-3-0

.PHONY: swagger-2-0
swagger-2-0: install-swag
	$(shell go env GOPATH)/bin/swag init \
		--generalInfo cmd/server/main.go \
		--dir . \
		--parseDependency \
		--parseInternal \
		--output docs/swagger \
		--generatedTime=false \
		--parseDepth 1 \
		--instanceName swagger \
		--parseVendor \
		--outputTypes go,json,yaml
	@make swagger-fix-refs

.PHONY: swagger-3-0
swagger-3-0: install-swag
	@echo "Converting Swagger 2.0 to OpenAPI 3.0..."
	@curl -X 'POST' \
		'https://converter.swagger.io/api/convert' \
		-H 'accept: application/json' \
		-H 'Content-Type: application/json' \
		-d @docs/swagger/swagger.json > docs/swagger/swagger-3-0.json
	@echo "Conversion complete. Output saved to docs/swagger/swagger-3-0.json"

.PHONY: swagger-fix-refs
swagger-fix-refs:
	@./scripts/fix_swagger_refs.sh

