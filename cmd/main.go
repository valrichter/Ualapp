package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valrichter/Ualapp/api"
	db "github.com/valrichter/Ualapp/db/sqlc"
	"github.com/valrichter/Ualapp/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:")
	}

	dbConnPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:")
	}

	store := db.NewPostgreSQLStore(dbConnPool)

	server, err := api.NewHTTPServer(store)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
	server.Start("localhost:8080")
}
