.PHONY: build run test clean fmt vet build-all

BINARY_NAME=gur

TARGETS = \
  "linux amd64" \
  "linux arm64" \
  "darwin amd64" \
  "darwin arm64" \
  "windows amd64" \
  "windows arm64"

build:
	go build -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Build for multiple OS/ARCH combos
build-all: clean
	mkdir -p dist
	@for target in $(TARGETS); do \
		set -- $$target; \
		GOOS=$$1 GOARCH=$$2 go build -C cmd/gur -o dist/$(BINARY_NAME)-$$1-$$2$(if $(findstring windows,$$1),.exe,) . ; \
		echo "Built: dist/$(BINARY_NAME)-$$1-$$2" ; \
	done