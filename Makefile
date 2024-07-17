postgres:
	docker run --name so-postgres -p 6666:5432 -e POSTGRES_PASSWORD=sopostgres -d postgres:16.3-alpine3.19

createdb:
	docker exec -it so-postgres createdb --username=postgres --owner=postgres go-so

dropdb:
	docker exec -it so-postgres dropdb --username=postgres go-so

migrateup:
	migrate -path db/migration -database "postgresql://postgres:sopostgres@localhost:6666/go-so?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:sopostgres@localhost:6666/go-so?sslmode=disable" -verbose down

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

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test installsqlc server mockgen