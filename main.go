package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/sunnyegg/go-so/api"
	db "github.com/sunnyegg/go-so/db/sqlc"
	"github.com/sunnyegg/go-so/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	serverAddress := config.ServerAddress
	dbSource := config.DBSource

	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}
	defer conn.Close(context.Background())

	queries := db.NewStore(conn)
	server := api.NewServer(queries)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
