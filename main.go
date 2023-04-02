package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mahmoud24598salah/MSM_Bank/api"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
	"github.com/mahmoud24598salah/MSM_Bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can nor read the config file", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to database", err)
	}

	store := db.NewStore(conn)

	server, err := api.Newserver(config, store)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	err = server.Start(config.SreverAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}

}
