package scripts

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valrichter/Ualapp/util"
)

func AddAccountNumbers() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:")
	}

	dbConnPool, err := pgxpool.New(context.Background(), config.DBSource+config.DBName+"?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db:")
	}

	accounts, err := dbConnPool.Query(context.Background(), "SELECT id,currency FROM accounts WHERE account_number IS NULL")
	if err != nil {
		panic(fmt.Sprint("cannot query accounts:", err))
	}

	for accounts.Next() {
		var id int32
		var currency string

		err = accounts.Scan(&id, &currency)
		if err != nil {
			panic(fmt.Sprint("cannot scan account:", err))
		}

		accountNumber, err := util.GenerateAccountNumber(id, currency)
		if err != nil {
			panic(fmt.Sprint("cannot generate account:", err))
		}

		_, err = dbConnPool.Exec(context.Background(), "UPDATE accounts SET account_number = $1 WHERE id = $2", accountNumber, id)
		if err != nil {
			panic(fmt.Sprint("cannot update account:", err))
		}
	}

	fmt.Println("finished adding account numbers")
}
