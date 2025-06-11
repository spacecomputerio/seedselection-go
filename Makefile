GOPACKAGES=$(shell go list ./...)
GOPACKAGES=$(shell go list ./...)
COVERPKG?=$(shell go list ./pkg/...)
BENCHPKG?=$(shell go list ./pkg/...)

deps:
	@go mod download

test: deps
	@go test -v $(GOPACKAGES)

race: deps
	@go test -race $(GOPACKAGES)

bench: deps
	@go test -bench=. -benchmem ${BENCHPKG} 

coverage: deps
	@go test -cover -coverprofile cover.out $(COVERPKG)
	@go tool cover -html=cover.out
	@rm -f cover.out

coverage-ci: deps
	@go test -cover $(COVERPKG)

fmt: deps
	@go fmt ./...

tidy: deps
	@go mod tidy

lint: deps
	@golangci-lint run ./...

build: deps
	@go build ./...

default: build