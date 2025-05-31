serve:
	CompileDaemon -command="./attendance-service" -color=true

run-migration:
	go run ./migrate/migrate.go

# Generate mocks
generate-mocks:
	mockery --config=.mockery.yaml
	@echo "Mocks generated in ./mocks directory"

# Run all tests
test:
	go test -v ./...