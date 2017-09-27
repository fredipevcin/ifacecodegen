PKG_LIST ?= $(shell go list ./... | grep -v '/vendor/'| grep -v '/examples/')
DIST = dist/ifacecodegen
test:
	go test -v -race -cover .

lint:
	golint $(PKG_LIST)

clean:
	rm -rf dist

build: $(DIST)

$(DIST):
	go build -ldflags="-s -w" -o dist/ifacecodegen ./cmd/ifacecodegen

run: $(DIST)
	@echo "Example 1"
	@./$(DIST) -source examples/interface.go -template examples/example1.tmpl -destination -
	@echo
	@echo "Example 2"
	@./$(DIST) -source examples/interface.go -template examples/example2.tmpl -destination -
	@echo

qa: lint test
