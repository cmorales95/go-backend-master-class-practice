package db

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"

	"github.com/cmorales95/go-backend-master-class/util"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database")
	}

	testQueries = New(testDB)

	m.Run()
}
