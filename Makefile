.PHONY: all
all: setup lint test

.PHONY: test
test: setup
	go test $$(go list ./... | grep -v acceptance)

.PHONY: testacc
testacc: setup
	go test ./acceptance

sources = $(shell find . -name '*.go' -not -path './vendor/*')
.PHONY: goimports
goimports: setup
	@goimports -w $(sources)

.PHONY: lint
lint: setup
	gometalinter $$(go list ./...) --enable=goimports --disable=gotype --disable=golint --disable=errcheck --vendor -t

.PHONY: errcheck
errcheck: setup
	gometalinter $$(go list ./...) --disable-all --enable=errcheck --vendor -t

.PHONY: install
install: setup
	go install

BIN_DIR := $(GOPATH)/bin
GOIMPORTS := $(BIN_DIR)/goimports
GOMETALINTER := $(BIN_DIR)/gometalinter
DEP := $(BIN_DIR)/dep
GOCOV := $(BIN_DIR)/gocov
GOCOV_HTML := $(BIN_DIR)/gocov-html

$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

$(DEP):
	go get -u github.com/golang/dep/cmd/dep

tools: $(GOIMPORTS) $(GOMETALINTER) $(DEP)

vendor: $(DEP)
	dep ensure

setup: tools vendor

updatedeps:
	dep ensure -update