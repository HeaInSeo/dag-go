.PHONY: test lint lint-fix fmt

# 로컬 bin에 golangci-lint가 있으면 우선 사용, 없으면 PATH에서 탐색
GOLANGCI_LINT ?= ./bin/golangci-lint

test:
	go test -v -race -cover ./...

bench:
	go test -bench=. -benchmem -benchtime=3s -run=^$ ./...

lint:
	$(GOLANGCI_LINT) run --config=.golangci.yml ./...

lint-fix:
	$(GOLANGCI_LINT) run --config=.golangci.yml --fix ./...

fmt:
	go fmt ./...