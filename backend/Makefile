.PHONY: test test-unit test-integration test-all build run clean

# Default target
all: build

# Build the application
build:
	go build -o bin/halooid cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run unit tests
test-unit:
	go test -v ./internal/... -short

# Run integration tests
test-integration:
	go test -v ./internal/taskodex/project -run TestProjectAPI
	go test -v ./internal/taskodex/task -run TestTaskAPI
	go test -v ./internal/taskodex/timeentry -run TestTimeEntryAPI
	go test -v ./internal/taskodex/comment -run TestCommentAPI
	go test -v ./internal/taskodex/fileattachment -run TestFileAttachmentAPI

# Run all tests
test-all: test-unit test-integration

# Run tests with default configuration
test: test-unit

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Create test database
create-test-db:
	psql -U postgres -c "DROP DATABASE IF EXISTS halooid_test;"
	psql -U postgres -c "CREATE DATABASE halooid_test;"
	migrate -path scripts/migrations -database "postgres://postgres:postgres@localhost:5432/halooid_test?sslmode=disable" up

# Drop test database
drop-test-db:
	psql -U postgres -c "DROP DATABASE IF EXISTS halooid_test;"
