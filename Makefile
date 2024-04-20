run:
	#PORT=8080 go run main.go
	DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable" PORT=8080 go run main.go

quality:
	go fmt ./...
	go vet ./...
	golangci-lint run