package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	dbSource = "postgresql://postgres:sopostgres@localhost:6666/go-so?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}
	defer conn.Close(context.Background())

	testQueries = New(conn)

	os.Exit(m.Run())
}
