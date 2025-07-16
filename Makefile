run:
	go run cmd/catalogmapperstage/main.go

build:
	cd cmd/catalogmapperstage && go build -o ../../catalogmapper

lint:
	golangci-lint run

createNewMigration:
	goose -dir migrations/postgres create first_migration sql

