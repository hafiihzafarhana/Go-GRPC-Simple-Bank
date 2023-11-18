postgres:
	docker run --name gomakemigrate -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it gomakemigrate createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it gomakemigrate dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

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

server:
	go run main.go

mockgen:
	mockgen -package mockdb -destination db/mock/store.go github.com/hafiihzafarhana/Go-GRPC-Simple-Bank/db/sqlc MockStore

goclean:
	gofmt -w ../.

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 startdb stopdb sqlcgen gotidy gotest server mockgen goclean