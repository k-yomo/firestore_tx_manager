.PHONY: test
test:
	FIRESTORE_EMULATOR_HOST=localhost:9090 go test -v -coverprofile cover.out ./...

.PHONY: up_store
up_store:
	docker-compose up -d
