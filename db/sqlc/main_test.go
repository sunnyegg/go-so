package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/sunnyegg/go-so/util"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	dbSource := config.DBSource
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}
	defer conn.Close(context.Background())

	testStore = NewStore(conn)

	os.Exit(m.Run())
}
