package db_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

var testStore db.Store

const testDBName = "test_ualapp"
const sslmode = "?sslmode=disable"

// TestMain sets up the database connection everytime the tests are run
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource+config.DBName+sslmode)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// create database for testing
	_, err = connPool.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s;", testDBName))
	if err != nil {
		log.Fatalf("Encountered an error creating database teardown (connPool) %v", err)
	}

	testConnPool, err := pgxpool.New(context.Background(), config.DBSource+testDBName+sslmode)
	if err != nil {
		log.Fatalf("Cannot connect to database %v", err)
		teardown(connPool)
	}

	conn, err := testConnPool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Cannot acquire connection %v", err)
		teardown(connPool)
	}
	connConfig := conn.Conn().Config()

	stdconn := stdlib.OpenDB(*connConfig)

	driver, err := postgres.WithInstance(stdconn, &postgres.Config{})
	if err != nil {
		log.Fatalf("Cannot create driver %v", err)
		teardown(connPool)
	}

	mig, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "../migrations"), config.DBDriver, driver)
	if err != nil {
		log.Fatalf("Cannot create migrate instance %v", err)
		teardown(connPool)
	}

	if err = mig.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Cannot migrate database %v", err)
		teardown(connPool)
	}

	testStore = db.NewPostgreSQLStore(testConnPool)

	code := m.Run()

	teardown(connPool)
	os.Exit(code)

	defer testConnPool.Close()
}

func teardown(conn *pgxpool.Pool) {
	_, err := conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s WITH (FORCE);", testDBName))
	if err != nil {
		log.Fatalf("failed to drop test database %v", err)
	}
	conn.Close()
}

// TODO: add tests for account.sql.go

// TODO: add tests for entries.sql.go
