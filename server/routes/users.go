package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	dbHelper "github.com/treavorgagne/twitter-clone/server/db"
	"github.com/treavorgagne/twitter-clone/server/models"
)

func GetUser(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()
	var user models.GetUserResponse

	row := conn.QueryRowContext(c.Request.Context(), "select * from users_stats where user_id = ?;", c.Param("user_id"))
	err = row.Scan(&user.User_id, &user.Username, &user.Created_At, &user.Total_following, &user.Total_followers)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Panic("ERROR_GETTING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()

	// Bind incoming JSON to the struct
	var req models.CreateUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	res, err := conn.ExecContext(c.Request.Context(), "insert into users (username) values (?);", req.Username)
	if err != nil {
		log.Panic("ERROR_CREATING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Fetch the new user ID
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("ERROR_FETCHING_INSERT_ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Println("users created username:", req.Username)
	c.JSON(http.StatusCreated, gin.H{"user_id": id})
}

func UpdateUser(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()

	var req models.UpdateUserRequest
	var user_id  = c.Param("user_id");
	// Bind incoming JSON to the struct
	if err = c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err = conn.ExecContext(c.Request.Context(), "update users SET username = ? where user_id = ?;", req.Username, user_id)
	if err != nil {
		log.Panic("ERROR_UPDATING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Println("users updated user_id:", user_id)
	c.Status(http.StatusNoContent)
}

func DeleteUser(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()

	user_id := c.Param("user_id")
	_, err = conn.ExecContext(c.Request.Context(), "delete from users where user_id = ?;", user_id)
	if err != nil {
		log.Panic("ERROR_DELETING_USER_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Println("users deleted user_id:", user_id)
	c.Status(http.StatusNoContent)
}