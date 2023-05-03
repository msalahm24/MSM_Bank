createdb:
	docker exec -it postgres createdb --username=root --owner=root msmBank

dropdb:
	docker exec -it postgres dropdb msmBank

migrateup:
	migrate -path db/migration -database "postgres://root:123@localhost:5432/msmBank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:123@localhost:5432/msmBank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

postgres:
	docker exec -it postgres psql -U root -d msmBank

server:
	go run main.go

mock:
	mockgen -package mockDB -destination db/mock/store.go github.com/mahmoud24598salah/MSM_Bank/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans:
	evans -r repl
	
.PHONY: createdb dropdb migrateup migratedown sqlc test server mock proto evans