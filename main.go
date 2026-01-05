package main

import (
	"database/sql"
	"log"
	db "github.com/techschool/simple-bank/db2/sqlc"
	"github.com/techschool/simple-bank/api"

	_ "github.com/lib/pq" // PostgreSQL driver,
	// _ is blank identifier to avoid import error if not used
	"github.com/techschool/simple-bank/utils"
)



func main() {
	config, err := utils.LoadConfig(".") // because config file is in same directory as main.go
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}