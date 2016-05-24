APP_NAME = kviator
VERSION = latest

all: package

clean:
	@echo "--> Cleaning build"
	@rm -rf ./build
	@rm -rf ./release

prepare:
	@echo "--> Preparing build"
	@mkdir -p build/bin/`go env GOOS`/`go env GOARCH`
	@mkdir -p build/test
	@mkdir -p build/doc
	@mkdir -p build/zip

format:
	@echo "--> Formatting source code"
	@go fmt ./...

deps:
	@echo "--> Getting dependencies"
	@go get -d -v ./...

test: prepare format deps
	@echo "--> Testing application"
	@go test -outputdir build/test ./...

build: prepare format deps
	@echo "--> Building local application"
	@go build -o build/bin/`go env GOOS`/`go env GOARCH`/${VERSION}/${APP_NAME} -ldflags "-X main.version=${VERSION}" -v .

package: build
	@echo "--> Packaging application"
	@zip -vj build/zip/${APP_NAME}-${VERSION}-`go env GOOS`-`go env GOARCH`.zip build/bin/`go env GOOS`/`go env GOARCH`/${VERSION}/${APP_NAME}
