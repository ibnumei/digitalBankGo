package main

import (
	"database/sql"
	"log"

	"github.com/ibnumei/digitalBankGo/util"

	"github.com/ibnumei/digitalBankGo/api"
	db "github.com/ibnumei/digitalBankGo/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("canot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to  db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot  start server", err)
	}
}
