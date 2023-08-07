package main

import (
	"database/sql"
	"log"

	"github.com/ibnumei/digitalBankGo/api"
	db "github.com/ibnumei/digitalBankGo/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:mysecretpassword@localhost:5432/master_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to  db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot  start server", err)
	}
}
