package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vbrenister/apicommon"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models

	apicommon.ServerConfig
}

func main() {
	log.Printf("Starting authentication services at port %s", webPort)

	conn := connecToDB()
	if conn == nil {
		log.Panic("Could not connect to the database")
	}

	app := &Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Panic(srv.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connecToDB() *sql.DB {
	dsn := os.Getenv("DNS")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready")
			counts++
		} else {
			log.Println("Postgres is ready")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Retrying in 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
