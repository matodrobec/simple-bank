name ?= init_schema  # Default value, can be overridden
type=gin
port=8090
DB_URL=postgresql://postgres:test@localhost:5432/bank?sslmode=disable

docker-compose:
	docker compose -p bank -f ./docker/docker-compose.yaml up -d --build

docker-compose-down:
	docker compose -p bank -f ./docker/docker-compose.yaml down

docker-compose-pbmodeler:
	docker compose -p bank -f ./docker/docker-compose-pgmodeler.yaml up -d --build

docker-compose-pbmodeler-down:
	docker compose -p bank -f ./docker/docker-compose-pgmodeler.yaml down

db-psql:
	docker compose -p bank -f ./docker/docker-compose.yaml exec postgres psql -U postgres -d bank

db-bash:
	docker compose -p bank -f ./docker/docker-compose.yaml exec postgres bash

db-logs-all:
	docker compose -p bank -f ./docker/docker-compose.yaml logs -f postgres

db-logs:
	docker compose -p bank -f ./docker/docker-compose.yaml logs -f -n 0 postgres

db-restart:
	docker compose -p bank -f ./docker/docker-compose.yaml restart postgres

db-create:
	docker compose -p bank -f ./docker/docker-compose.yaml exec postgres createdb --username=postgres --owner=postgres --template=template0 bank

db-drop:
	docker compose -p bank -f ./docker/docker-compose.yaml exec postgres dropdb --username=postgres bank

migrate: migrate-up sqlc-gen mock

migrate-up:
	migrate --path db/migration --database "$(DB_URL)" --verbose up

migrate-down:
	migrate --path db/migration --database "$(DB_URL)" --verbose down

migrate-down-last:
	migrate --path db/migration --database "$(DB_URL)" --verbose down 1

migrate-create:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc: sqlc-gen mock

sqlc-gen:
	sqlc generate

test:
	go test -v -cover -short ./db/... ./api/... ./util/... ./mail/...

server:
	go run main.go $(type)

mock:
	mockgen --destination db/mock/store.go -package mockdb  github.com/matodrobec/simplebank/db/sqlc Store
	mockgen --destination worker/mock/distributor.go -package mockwk github.com/matodrobec/simplebank/worker TaskDistributor

proto: proto-run statik-swagger

proto-run:
	rm -f pb/*.go
	rm -f doc/swagger/*.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out ./pb \
		--grpc-gateway_opt paths=source_relative \
		--openapiv2_out ./doc/swagger \
		--openapiv2_opt allow_merge=true,json_names_for_fields=false,merge_file_name=simple_bank  \
    proto/*.proto

proto-evans:
	evans --host localhost --port $(port) -r repl

statik-swagger:
	statik -src=./doc/swagger -dest=./doc -f -ns swagger


.PHONY: docker-compose docker-compose-down db-restart db-create db-drop migrate-up migrate-down sqlc-gen test server mock docker-compose-pbmodeler migrate-create proto proto-evans statik-swagger