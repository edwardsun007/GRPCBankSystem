package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver,
	// _ is blank identifier to avoid import error if not used
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// this is a global variable that will be used by all unit tests
var testQueries *Queries
var testDB *sql.DB // Store the database connection for raw SQL queries

// entrance point for all unit tests
func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer conn.Close()

	// Test the connection
	if err = conn.Ping(); err != nil {
		log.Fatal("Cannot ping database:", err)
	}

	testQueries = New(conn) // this calls the New function defined in db.go file
	testDB = conn           // Store connection for raw SQL queries

	os.Exit(m.Run())
}
