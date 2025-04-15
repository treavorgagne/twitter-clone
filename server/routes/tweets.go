package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/treavorgagne/twitter-clone/server/models"
)

func GetTweet(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tweet_id := c.Param("tweet_id")
		var tweet models.GetTweetsResponse

		row := db.QueryRow("select * from tweets_stats where tweet_id = ?;", tweet_id)
		err := row.Scan(&tweet.User_id, &tweet.Body, &tweet.User_id, &tweet.Created_At, &tweet.Tweet_total_likes);
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
			return
		} else if err != nil {
			log.Printf("ERROR_SCANNING_TWEET_ROW: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tweet data"})
			return
		}

		c.JSON(http.StatusOK, tweet)
	}
}

func GetTweets(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")

		rows, err := db.Query("select * from tweets_stats where user_id = ?;", user_id)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No Tweets found"})
			return
		} else if err != nil {
			log.Panic("ERROR_GETTING_TWEETS_DATA: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		var tweets []models.GetTweetsResponse
		for rows.Next() {
			var tweet models.GetTweetsResponse
			if err := rows.Scan(&tweet.User_id, &tweet.Body, &tweet.User_id, &tweet.Created_At, &tweet.Tweet_total_likes); err != nil {
				log.Printf("ERROR_SCANNING_TWEET_ROW: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tweet data"})
				return
			}
			tweets = append(tweets, tweet)
		}

		c.JSON(http.StatusOK, tweets)
	}
}

func CreateTweet(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateTweetRequest
		user_id := c.Param("user_id")
		// Bind incoming JSON to the struct
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := db.Exec("insert into tweets (body, user_id) values (?, ?);", req.Body, user_id)
		if err != nil {
			log.Panic("ERROR_CREATING_TWEET: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		log.Println("tweet created by user_id:", user_id)
		c.Status(http.StatusCreated)
	}
}

func UpdateTweet(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UpdateTweetRequest
		user_id := c.Param("user_id")
		tweet_id := c.Param("tweet_id")
		// Bind incoming JSON to the struct
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		_, err := db.Exec("update tweets SET body = ? where tweet_id = ? and user_id = ?;", req.Body, tweet_id, user_id)
		if err != nil {
			log.Panic("ERROR_UPDATING_TWEET: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		log.Println("updated tweet tweet_id:", tweet_id)
		c.Status(http.StatusNoContent)
	}
}

func DeleteTweet(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		tweet_id := c.Param("tweet_id")
		_, err := db.Exec("delete from tweets where user_id = ? and tweet_id = ?;", user_id, tweet_id)
		if err != nil {
			log.Panic("ERROR_DELETING_TWEET_DATA: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		log.Println("tweet deleted user_id:", user_id)
		c.Status(http.StatusNoContent)
	}
}

func LikeTweet(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		tweet_id := c.Param("tweet_id")
		_, err := db.Exec("insert into likes_tweets (tweet_id, user_id) values (?,?);", tweet_id, user_id)
		if err != nil {
			log.Panic("ERROR_LIKING_TWEET_DATA: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		log.Printf("tweet (%s) liked by user_id: (%s)", tweet_id, user_id)
		c.Status(http.StatusNoContent)
	}
}

func UnLikeTweet(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		tweet_id := c.Param("tweet_id")
		_, err := db.Exec("delete from likes_tweets where user_id = ? and tweet_id = ?;", user_id, tweet_id)
		if err != nil {
			log.Panic("ERROR_DELETING_LIKE_DATA: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		log.Printf("tweet (%s) unliked by user_id: (%s)", tweet_id, user_id)
		c.Status(http.StatusNoContent)
	}
}
