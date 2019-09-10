.PHONY: test
test:
	FIRESTORE_EMULATOR_HOST=firestore:9090 go test -v -coverprofile cover.out ./...

.PHONY: up_store
up_store:
	docker-compose up -d
