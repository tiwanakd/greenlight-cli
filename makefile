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
	@echo 'Formatting .go files'
	go fmt ./...
	@echo 'Tidying module dependencies'
	go mod tidy

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
