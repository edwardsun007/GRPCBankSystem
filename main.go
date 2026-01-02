package main

import (
	"database/sql"
	"log"
	db "github.com/techschool/simple-bank/db2/sqlc"
	"github.com/techschool/simple-bank/api"

	_ "github.com/lib/pq" // PostgreSQL driver,
	// _ is blank identifier to avoid import error if not used
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080" // TODO: change to env variable
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}