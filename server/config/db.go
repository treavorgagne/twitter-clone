package config

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func ConfigDB() *sql.DB {
    dbAddress := os.Getenv("DBADDRESS")
    dbUser := os.Getenv("DBUSER")
	log.Println(dbAddress, dbUser)

    cfg := mysql.Config{
        User:   dbUser,
        Passwd: os.Getenv("DBPASS"),
        Net:    "tcp",
        Addr:   dbAddress + ":3306",
        DBName: "twitter",
    }

    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatalf("DB_CONNECTION_FAILED: unable to open database connection: %v", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("DB_CONNECTION_FAILED: unable to ping database at %s: %v", cfg.Addr, err)
    }

    log.Println("Database connection established to", cfg.Addr)
    return db
}


func GetDBConn(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := db.Conn(c) // new scoped connection
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "DB connection failed"})
			return
		}
		c.Set("conn", conn)
		defer conn.Close()
		c.Next()
	}
}