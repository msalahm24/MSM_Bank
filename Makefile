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

.PHONY: createdb dropdb migrateup migratedown sqlc