run:
	DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable" PORT=8080 go run main.go

quality:
	go fmt ./...
	go vet ./...
	golangci-lint run

test-unit:
	go test -v ./... -tags=unit

test-integration:
	docker-compose -f docker-compose.it-test.yaml down && \
	docker-compose -f docker-compose.it-test.yaml up --build --force-recreate --abort-on-container-exit --exit-code-from app_integration_tests

test-cover:
	go test -tags=unit -cover ./...

test-cover-html:
	go test -tags=unit -coverprofile=c.out ./...
	go tool cover -html=c.out
	rm c.out

docker-build:
	docker build -t ktaxes .

docker-run:
	docker run -p 8080:8080 -e DATABASE_URL="host=postgres port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable" ktaxes

swagger:
	swag init