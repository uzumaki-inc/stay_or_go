.PHONY: air lint lintFix cover

# airを実行
air:
	go run github.com/air-verse/air

# golangci-lintを実行
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run
lintFix:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix

# Run tests with coverage and open HTML report
cover:
	mkdir -p .gocache
	GOCACHE=$(PWD)/.gocache go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@if command -v open >/dev/null 2>&1; then \
		open coverage.html; \
	elif command -v xdg-open >/dev/null 2>&1; then \
		xdg-open coverage.html; \
	else \
		echo "Coverage report saved to $(PWD)/coverage.html"; \
	fi
