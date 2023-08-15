package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	// "github.com/spf13/viper"

	"github.com/ibnumei/digitalBankGo/api"
	db "github.com/ibnumei/digitalBankGo/db/sqlc"
	"github.com/ibnumei/digitalBankGo/util"
)

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgresql://root:mysecretpassword@localhost:5432/master_bank?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// )

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
