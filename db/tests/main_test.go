package db_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

var testStore db.Store

// TestMain sets up the database connection everytime the tests are run
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	connPoll, err := pgxpool.New(context.Background(), config.DBSourceTest)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	testStore = db.NewPostgreSQLStore(connPoll)
	defer connPoll.Close()
	os.Exit(m.Run())
}

// TODO: add tests for account.sql.go

// TODO: add tests for entries.sql.go
