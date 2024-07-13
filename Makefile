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

sqlc:
	sqlc generate
	
test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test