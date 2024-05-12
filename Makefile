DB_URL=postgresql://root:secret@localhost:5455/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 -p 5455:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
mysql:
	docker run --name mysql8 -p 3306:3306  -e MYSQL_ROOT_PASSWORD=secret -d mysql:8
	
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose  down

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose  down 1

db_docs:
	dbdocs build doc/db.dbml 

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
	
sqlc:
	sqlc generate
run:
	go run main.go

test: 
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb  -destination=db/mock/store.go github/bekeeeee/simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	 proto/*.proto

evans:
	evans --host 0.0.0.0 --port 9090 -r repl

redis:
	docker run --name redis -p 6339:6379 -d redis:7-alpine
start_redis:
	docker start redis
stop_redis:
	docker stop redis

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1 migrateupN migratedownN db_docs db_schema proto evans redis start_redis stop_redis

