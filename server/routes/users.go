package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/treavorgagne/twitter-clone/server/models"
)

func GetUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.GetUserResponse

		row := db.QueryRow("select * from users_stats where user_id = ?;", c.Param("user_id"))
		err := row.Scan(&user.User_id, &user.Username, &user.Created_At, &user.Total_following, &user.Total_followers)
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
}

func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateUserRequest

		// Bind incoming JSON to the struct
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := db.Exec("insert into users (username) values (?);", req.Username)
		if err != nil {
			log.Panic("ERROR_GETTING_USER_DATA: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		log.Println("users created username:", req.Username)
		c.Status(http.StatusCreated)
	}
}

func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UpdateUserRequest
		var user_id  = c.Param("user_id");

		// Bind incoming JSON to the struct
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := db.Exec("update users SET username = ? where user_id = ?;", req.Username, user_id)
		if err != nil {
			log.Panic("ERROR_UPDATING_USER_DATA: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		log.Println("users updated user_id:", user_id)
		c.Status(http.StatusNoContent)
	}
}

func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		_, err := db.Exec("delete from users where user_id = ?;", user_id)
		if err != nil {
			log.Panic("ERROR_DELETING_USER_DATA: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		log.Println("users deleted user_id:", user_id)
		c.Status(http.StatusNoContent)
	}
}