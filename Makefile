.PHONY: createdb dropdb postgres migrateup migrateup-last migratedown migratedown-last server mock random-symmetric-key db-docs db-schema proto evans

DB_URL = postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

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
	migrate --path db/migratpon -database $(DB_URL) -verbose up

migrateup-last:
	migrate --path db/migration -database $(DB_URL) -verbose up 1

migratedown:
	migrate --path db/migration -database $(DB_URL) -verbose down

migratedown-last:
	migrate --path db/migration -database $(DB_URL) -verbose down 1
sqlc:
	# https://docs.sqlc.dev/en/latest/overview/install.html#ubuntu
	sqlc generate
	make mock

test:
	go test -v -cover ./... -count=1

server:
	@go run main.go

mock:
	mockgen -destination=db/mocks/store.go -package=mocks github.com/cmorales95/go-backend-master-class/db/sqlc Store

random-symmetric-key:
	openssl rand -hex 64 | head -c 32

db-docs:
	# installation: npm install -g dbdocs  | dbdocs login
	dbdocs build docs/db.dbml

db-schema:
	# installation: npm install -g @dbml/cli
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

proto:
	@rm -f pb/*.go
	@rm -f docs/swagger/*.swagger.json
	@rm -f docs/statik/statik.go
	@protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
		--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
		proto/*.proto
	@statik -src=./docs/swagger -dest=./docs

evans:
	evans -r repl