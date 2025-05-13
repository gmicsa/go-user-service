version="0.0.1"

.PHONY: build run imports lint test

default: imports lint test build

build:
	docker build --build-arg APP_VERSION=$(version) -t user-service:$(version) .

run:
	docker run -p 8080:8080 -p 8081:8081 -p 8082:8082 --rm user-service:$(version)

imports:
	git ls-files '*.go' ':!vendor/*' | xargs goimports -w

lint:
	golangci-lint run ./...

test:
	go test ./...