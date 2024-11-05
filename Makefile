test:
	bash -c "set -m; bash '$(CURDIR)/scripts/test.sh'"

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 run

mocks:
	go run github.com/vektra/mockery/v2@v2.46.3

format:
	go mod tidy
	go fmt ./...
	go run github.com/daixiang0/gci@latest write \
		--skip-generated \
		-s standard -s default \
		-s "prefix(github.com/a-novel/golib)" \
		-s "prefix(buf.build/gen/go/a-novel)" \
		-s "prefix(github.com/a-novel/uservice-passkeys)" \
		.
	go run mvdan.cc/gofumpt@latest -l -w .

run:
	bash -c "set -m; bash '$(CURDIR)/scripts/run.sh'"

.PHONY: run test lint format
