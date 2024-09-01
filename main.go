package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sunnyegg/go-so/api"
	"github.com/sunnyegg/go-so/channel"
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

	// create connection pool
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}
	defer conn.Close()

	queries := db.NewStore(conn)
	server, err := api.NewServer(config, queries)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	// create channels
	channel.NewChannel(channel.ChannelWebsocket).Create()
	channel.NewChannel(channel.ChannelBlacklist).Create()
	channel.NewChannel(channel.ChannelEventsub).Create()
	channel.NewChannel(channel.ChannelGeneral).Create()

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
