GOPATH:=$(shell go env GOPATH)

.PHONY: init
	@go mod tidy

snapshots:
	@goreleaser check
	@goreleaser release --snapshot --skip-publish --rm-dist

