package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func connectDB(username, password, host, database string, port int) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, database)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}

func main() {
	db := connectDB("postgres", "password", "localhost", "postgres", 5432)
	repo := newPostgresRepo(db)
	service := newPatientsService(repo)
	httpTransport := newHttpTransport(service)

	routes := buildRoutes(httpTransport)

	err := http.ListenAndServe(":8000", routes)
	log.Println("Some error occured while listening to port 8000:", err)
}
