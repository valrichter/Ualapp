package main

import (
	"context"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valrichter/Ualapp/api"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/scripts"
	"github.com/valrichter/Ualapp/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:")
	}

	dbConnPool, err := pgxpool.New(context.Background(), config.DBSource+config.DBName+"?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db:")
	}

	runDBMigration(config.MigrationURL, config.DBSource+config.DBName+"?sslmode=disable")

	store := db.NewPostgreSQLStore(dbConnPool)

	server, err := api.NewHTTPServer(store)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
	server.Start("0.0.0.0:8080")
	scripts.AddAccountNumbers()
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create migration: ", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("cannot run migration: ", err)
	}

	log.Println("db migrated")
}
