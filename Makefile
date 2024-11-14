.PHONY: air lint

# airを実行
air:
	go run github.com/air-verse/air

# golangci-lintを実行
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run
lintFix:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix
