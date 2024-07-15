package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/sunnyegg/go-so/api"
	db "github.com/sunnyegg/go-so/sqlc"
)

const (
	dbSource      = "postgresql://postgres:sopostgres@localhost:6666/go-so?sslmode=disable"
	serverAddress = "0.0.0.0:9000"
)

func main() {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}
	defer conn.Close(context.Background())

	queries := db.New(conn)
	server := api.NewServer(queries)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
