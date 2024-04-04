dev:
	go build -o out/dbt-lsp .

prod:
	go build -ldflags "-s -w" -o out/prod/dbt-lsp

test:
	go test . -v
