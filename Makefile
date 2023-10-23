postgres:
	docker run --name gomakemigrate -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it gomakemigrate createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it gomakemigrate dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

startdb:
	docker start gomakemigrate

stopdb:
	docker stop gomakemigrate

sqlcgen:
	sqlc generate

gotidy:
	go mod tidy

gotest:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown startdb stopdb sqlcgen gotidy gotest