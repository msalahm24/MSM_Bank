package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:123@localhost:5432/msmBank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB
func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("can not connect to database", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
