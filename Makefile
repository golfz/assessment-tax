run:
	DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable" PORT=8080 go run main.go

quality:
	go fmt ./...
	go vet ./...
	golangci-lint run

test:
	go test -v ./...

integration-test:
	go test -v ./... -tags=integration

cover:
	go test -coverprofile=c.out ./...
	go tool cover -html=c.out
	rm c.out