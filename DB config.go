package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=devconnect sslmode=disable"
	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connected!")
}
