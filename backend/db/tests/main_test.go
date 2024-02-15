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

var testQuery db.Store

// TestMain sets up the database connection everytime the tests are run
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	// TODO: change this to use in test database
	connPoll, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	testQuery = db.NewPostgreSQLStore(connPoll)
	defer connPoll.Close()
	os.Exit(m.Run())
}
