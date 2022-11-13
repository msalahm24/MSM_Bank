package main

import (
	"database/sql"
	"log"

	
	_ "github.com/lib/pq"
	"github.com/mahmoud24598salah/MSM_Bank/api"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:123@localhost:5432/msmBank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can not connect to database", err)
	}

	store := db.NewStore(conn)

	server :=api.Newserver(store)

	err= server.Start(serverAddress)
	if err != nil{
		log.Fatal("Cannot start server",err)
	}

}
