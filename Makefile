GO_FILES = $(shell find . -name '*.go' -not -path "./vendor/*")
BINARY_NAME = app

.PHONY: run all build proto lint clean create-repository

run:
	@echo "🚀 Rodando o projeto..."
	go run cmd/main.go

all: build

build:
	@echo "🔨 Compiling..."
	go build -o $(BINARY_NAME) cmd/main.go

lint:
	@echo "🧐 Checking formatting and linting..."
	gofmt -w $(GO_FILES)
	go vet ./...

proto:
	@echo "🔨 Copying proto files..."
	cp ../ad-apis/gen/go/go_notifications.pb.go ./internal/infrastructure/grpc/notifications
	cp ../ad-apis/gen/go/go_notifications_grpc.pb.go ./internal/infrastructure/grpc/notifications

clean:
	@echo "🧹 Cleaning..."
	rm -f $(BINARY_NAME)

create-repository:
	@echo "🏗️ Creating ECR repository..."
	aws ecr create-repository --repository-name go-notifications-dev
