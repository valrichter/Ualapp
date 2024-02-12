package db_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	db "github.com/valrichter/Ualapp/db/sqlc"
)

var testQuery *db.Queries

// TestMain sets up the database connection everytime the tests are run
func TestMain(m *testing.M) {
	var dbURL string = "postgresql://root:secret@localhost:5432/ualapp?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}
	defer conn.Close(context.Background())

	testQuery = db.New(conn)
	os.Exit(m.Run())
}
