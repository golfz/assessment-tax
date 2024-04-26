quality:
	go fmt ./...
	go vet ./...
	golangci-lint run

test-unit:
	go test ./... -tags=unit

test-unit-v:
	go test -v ./... -tags=unit

test-cover:
	go test -tags=unit -cover ./...

test-cover-html:
	go test -tags=unit -coverprofile=c.out ./...
	go tool cover -html=c.out
	rm c.out

test-integration:
	docker-compose -f docker-compose.it-test.yaml down && \
	docker-compose -f docker-compose.it-test.yaml up --build --force-recreate --abort-on-container-exit --exit-code-from it-test-goapp

test-e2e:
	docker-compose -f docker-compose.e2e-test.yaml down && \
	docker-compose -f docker-compose.e2e-test.yaml up --build --force-recreate --abort-on-container-exit --exit-code-from e2e-postman

docker-build:
	docker build -t ktaxes .

docker-run-image:
	docker run -p 8080:8080 -e DATABASE_URL="host=postgres port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable" ktaxes

run-db:
	docker-compose -f docker-compose.yaml down && \
	docker-compose -f docker-compose.yaml up

down-db:
	docker-compose -f docker-compose.yaml down

run-app:
	DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable" PORT=8080 ADMIN_USERNAME=adminTax ADMIN_PASSWORD=admin! go run main.go

run-local:
	docker-compose -f docker-compose.local.yaml down && \
    docker-compose -f docker-compose.local.yaml up --build --force-recreate --abort-on-container-exit

down-local:
	docker-compose -f docker-compose.local.yaml down

swagger:
	swag init