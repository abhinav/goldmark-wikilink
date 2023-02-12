# "go install"-ed binaries will be placed here during development.
export GOBIN ?= $(shell pwd)/bin

GO_FILES = $(shell find . \
	   -path '*/.*' -prune -o \
	   '(' -type f -a -name '*.go' ')' -print)

GOLINT = $(GOBIN)/golint
STATICCHECK = $(GOBIN)/staticcheck
TOOLS = $(GOLINT) $(STATICCHECK)

.PHONY: all
all: build lint test

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: lint
lint: gofmt golint staticcheck

.PHONY: gofmt
gofmt:
	$(eval FMT_LOG := $(shell mktemp -t gofmt.XXXXX))
	@gofmt -e -s -l $(GO_FILES) > $(FMT_LOG) || true
	@[ ! -s "$(FMT_LOG)" ] || \
		(echo "gofmt failed. Please reformat the following files:" | \
		cat - $(FMT_LOG) && false)

.PHONY: golint
golint: $(GOLINT)
	$(GOLINT) ./...

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	$(eval STATICCHECK_LOG := $(shell mktemp -t staticcheck.XXXXX))
	$(STATICCHECK) ./... | grep -v SA1019 > $(STATICCHECK_LOG) || true
	@[ ! -s "$(STATICCHECK_LOG)" ] || \
		 (echo "static failed:" | \
		 cat - $(STATICCHECK_LOG) && false)

tools: $(GOLINT) $(STATICCHECK)

$(GOLINT): tools/go.mod
	cd tools && go install golang.org/x/lint/golint

$(STATICCHECK): tools/go.mod
	cd tools && go install honnef.co/go/tools/cmd/staticcheck
