.PHONY: createdb dropdb postgres migrateup migrateup-last migratedown migratedown-last server mock

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

infra:
	docker compose up postgres -d

infra-down:
	docker compose down

create-migration-file:
	migrate create -ext sql -dir db/migration -seq $(name) # call it using name=

migrateup:
	# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup-last:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown-last:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	# https://docs.sqlc.dev/en/latest/overview/install.html#ubuntu
	sqlc generate
	make mock

test:
	go test -v -cover ./... -count=1

server:
	go run main.go

mock:
	mockgen -destination=db/mocks/store.go -package=mocks github.com/cmorales95/go-backend-master-class/db/sqlc Store

random-symmetric-key:
	openssl rand -hex 64 | head -c 32