postgres:
	docker run --name so-postgres -p 6666:5432 -e POSTGRES_PASSWORD=sopostgres -d postgres:16.3-alpine3.19

createdb:
	docker exec -it so-postgres createdb --username=postgres --owner=postgres go-so

dropdb:
	docker exec -it so-postgres dropdb --username=postgres go-so

migrateup:
	migrate -path db/migration -database "postgresql://postgres:sopostgres@localhost:6666/go-so?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:sopostgres@localhost:6666/go-so?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://postgres:sopostgres@localhost:6666/go-so?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:sopostgres@localhost:6666/go-so?sslmode=disable" -verbose down 1

installsqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0

sqlc:
	sqlc generate
	
test:
	go test -v -cover ./...

server:
	go run main.go

mockgen:
	go run github.com/golang/mock/mockgen@v1.6.0 -destination db/mock/store.go -package mockdb github.com/sunnyegg/go-so/db/sqlc Store

createmigration:
	migrate create -ext sql -dir db/migration -seq $(seq)

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test installsqlc server mockgen migrateup1 migratedown1 createmigration