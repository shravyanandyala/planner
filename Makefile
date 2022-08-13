GO_SRC=$(shell find . -path ./.build -prune -false -o -name \*.go)
VERSION=$(shell git describe --tags --always || git rev-parse HEAD)
VERSION_FULL=$(if $(shell git status --porcelain --untracked-files=no),$(VERSION)-dirty,$(VERSION))

BUILD_TAGS = osusergo netgo static_build

build_planner = go build -tags "$(BUILD_TAGS)" -buildmode=pie -ldflags "-X main.version=$(VERSION_FULL) -extldflags '-static'" -o planner ./

planner: $(GO_SRC) go.mod go.sum
	$(call build_planner)

bin: planner

lint: ./.golangcilint.yaml
	golangci-lint --config ./.golangcilint.yaml run ./...

test: $(GO_SRC)
	go test -v -race -cover -coverpkg ./... -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: check
check: test lint

.PHONY: mod
mod:
	go get -u
	go mod tidy

.PHONY: clean
clean:
	rm ./planner
	rm ./coverage.txt