package db_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/valrichter/Ualapp/db/sqlc"
)

var testQuery db.Store

// TODO: put this in a config file
// TestMain sets up the database connection everytime the tests are run
func TestMain(m *testing.M) {
	const dbURL string = "postgresql://root:secret@localhost:5432/ualapp?sslmode=disable"

	connPoll, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	testQuery = db.NewPostgreSQLStore(connPoll)
	defer connPoll.Close()
	os.Exit(m.Run())
}
