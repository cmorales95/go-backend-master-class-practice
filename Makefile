.PHONY: createdb dropdb postgres migrateup migratedown server mock

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	# https://docs.sqlc.dev/en/latest/overview/install.html#ubuntu
	sqlc generate

test:
	go test -v -cover ./... -count=1

server:
	go run main.go

mock:
	mockgen -destination=db/mocks/store.go -package=mocks github.com/cmorales95/go-backend-master-class/db/sqlc Store