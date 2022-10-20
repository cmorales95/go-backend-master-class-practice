package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/cmorales95/go-backend-master-class/api"
	db "github.com/cmorales95/go-backend-master-class/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	address  = "localhost:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to the database")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
