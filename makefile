# ========================================================================================= #
# HELPERS
# ========================================================================================= #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ========================================================================================= #
# Delevelopment 
# ========================================================================================= #

## run/client: run the cli client
.PHONY: run/cli
run/cli:
	go run .

# ========================================================================================= #
# QUALITY CONTROL 
# ========================================================================================= #

## tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor

## audit: run quality control check
.PHONY: audit
audit: 
	@echo 'Checking module dependencies'
	go mod tidy -diff
	go mod verify
	@echo 'vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# ========================================================================================= #
# BUILD 
# ========================================================================================= #

## build/cli: Build the cli application
.PHONY: build/cli
build/cli:
	@echo 'Building cli application...'
	go build -ldflags='-s' -o=./bin/greenlight .
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/greenlight .
	GOOS=windows GOARCH=amd64 go build -ldflags='-s' -o=./bin/windows_amd64/greenlight.exe .
	GOOS=darwin GOARCH=amd64 go build -ldflags='-s' -o=./bin/darwin_amd64/greenlight .
