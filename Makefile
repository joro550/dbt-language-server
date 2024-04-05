dev:
	go build -o out/dbt-lsp .

windows:
	go build -o out/dbt-lsp.exe .

prod:
	go build -ldflags "-s -w" -o out/prod/dbt-lsp

test:
	go test . -v
