postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecretpassword -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root master_bank

dropdb:
	docker exec -it postgres12 dropdb master_bank

migrateup:
	migrate -path db/migration -database  "postgresql://root:mysecretpassword@localhost:5432/master_bank?sslmode=disable" -verbose up

migrateup1Version:
	migrate -path db/migration -database  "postgresql://root:mysecretpassword@localhost:5432/master_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database  "postgresql://root:mysecretpassword@localhost:5432/master_bank?sslmode=disable" -verbose down

migratedown1Version:
	migrate -path db/migration -database  "postgresql://root:mysecretpassword@localhost:5432/master_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ibnumei/digitalBankGo/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1Version migratedown1Version