package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/mahmoud24598salah/MSM_Bank/util"
)



var testQueries *Queries
var testDB *sql.DB
func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil{
		log.Fatal("Can not load config file",err)
	}
	testDB, err = sql.Open(config.DBDriver,config.DBSource)
	if err != nil {
		log.Fatal("can not connect to database", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
