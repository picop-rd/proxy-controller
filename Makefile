.PHONY:docker-compose-build
docker-compose-build: test
	docker compose build --ssh default

.PHONY:docker-compose-run
docker-compose-run: docker-compose-build
	docker compose up

.PHONY:test
test: vet fmt
	go test ./...

.PHONY:test-with-coverage
test-with-coverage: vet fmt
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

.PHONY:vet
vet:
	go vet ./...

.PHONY:fmt
fmt:
	go fmt ./...
