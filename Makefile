postgres:
	docker run --name postgres_container -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine3.21

createdb:
	docker exec -it postgres_container createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres_container dropdb simple_bank

migrateup:
	migrate -path db2/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db2/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db2/mock/store.go github.com/techschool/simple-bank/db2/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock

# WHY WE USE .PHONY:
# With .PHONY:
# Make always runs the commands, even if a file with that name exists
# Make won't check for a file with that name
# The targets are treated as commands/actions, not file outputs
