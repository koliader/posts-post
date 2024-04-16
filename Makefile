createdb:
	 docker exec -it pg12 createdb --username=root --owner=root postsPosts
dropdb:
	 docker exec -it pg12 dropdb postsPosts
migrateup:
	migrate -path internal/db/migration -database "postgresql://root:secret@127.0.0.1:5432/postsPosts?sslmode=disable" -verbose up
migratedown:
	migrate -path internal/db/migration -database "postgresql://root:secret@127.0.0.1:5432/postsPosts?sslmode=disable" -verbose down
sqlc:
	sqlc generate
protoc:
	rm -f ./internal/pb/*.go
	mkdir -p ./internal/pb
	protoc -I ./proto \
	--go_out ./internal/pb --go_opt paths=source_relative \
	--go-grpc_out ./internal/pb --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./internal/pb --grpc-gateway_opt paths=source_relative \
	proto/*.proto
test:
	go test -v -cover ./...
evans:
	 evans --host localhost --port 8081 -r repl
server:
	go run cmd/main.go
.PHONY: protoc
