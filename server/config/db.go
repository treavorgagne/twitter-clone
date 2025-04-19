package config

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func ConfigDB() (*sql.DB) {
	cfg := mysql.Config{
        User:   os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net:    "tcp",
        Addr:   os.Getenv("DBADDRESS")+":"+os.Getenv("DBPORT"),
        DBName: "twitter",
    }
    var db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
		log.Fatal("DB_CONNECTION_FAILED", "unable to get database connection")
	}
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