run: protoc
	go run cmd/catalogmapperstage/main.go

build: protoc
	cd cmd/catalogmapperstage && go build -o ../../catalogmapper

lint: protoc
	golangci-lint run

createNewMigration:
	goose -dir migrations/postgresql create first_migration sql

protoc:
	@echo "Generating protobuf files for service..."
	protoc --proto_path=api/ \
		   --proto_path=pkg/third_party/proto/googleapis \
	       --go_out=pkg/pb/ --go_opt=paths=source_relative \
	       --go-grpc_out=pkg/pb/ --go-grpc_opt=paths=source_relative \
	       api/exchangerateservice/*.proto