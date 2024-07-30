VERSION=0.1.0

BINARY_NAME=tt
MAIN_PACKAGE_PATH=.
PACKAGE=github.com/schretzi/tt

VERSION`git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')`
COMMIT_HASH=`git rev-parse --short HEAD`
BUILD_TIMESTAMP=`date '+%Y-%m-%dT%H:%M:%S'`

LDFLAGS=-ldflags "-X ${PACKAGE}/tt.Version=${VERSION} -X ${PACKAGE}/tt.CommitHash=${COMMIT_HASH} -X ${PACKAGE}/tt.BuildTimestamp=${BUILD_TIMESTAMP}"


# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

# ==================================================================================== #
# TEST & BUILD
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...


## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
    # Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build ${LDFLAGS} -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
