package db

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConfigDB() (*sql.DB) {
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := mysql.Config{
        User:   os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net:    "tcp",
        Addr:   "127.0.0.1:3308",
        DBName: "twitter",
    }
    var db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
		log.Fatal("DB_CONNECTION_FAILED", "unable to get database connection")
	}
    return db
}

func GetDBConn(db *sql.DB) (*sql.Conn, error) {
	return db.Conn(context.Background())
}