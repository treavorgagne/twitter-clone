package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	dbHelper "github.com/treavorgagne/twitter-clone/server/db"
	"github.com/treavorgagne/twitter-clone/server/models"
)

func GetTweet(c *gin.Context, db *sql.DB, rdb *redis.Client) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()

	tweet_id := c.Param("tweet_id")
	var tweet models.GetTweetsResponse

	row := conn.QueryRowContext(c.Request.Context(), "select * from tweets_stats where tweet_id = ?;", tweet_id)
	err = row.Scan(&tweet.Tweet_id, &tweet.Body, &tweet.User_id, &tweet.Created_At, &tweet.Tweet_total_likes);
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		return
	} else if err != nil {
		log.Printf("ERROR_SCANNING_TWEET_ROW: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tweet data"})
		return
	}

	// 🔥 Store into Redis cache
	// 🧹 Marshal struct to JSON before caching
	tweetJSON, err := json.Marshal(tweet)
	if err != nil {
		log.Println("Failed to marshal user:", err)
	} else {
		cacheKey := c.Request.URL.Path
		err := rdb.Set(c, cacheKey, tweetJSON, 5*time.Minute).Err()
		if err != nil {
			log.Println("Redis cache SET error:", err)
		}
	}

	c.JSON(http.StatusOK, tweet)
}

func GetTweets(c *gin.Context, db *sql.DB, rdb *redis.Client) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()
	user_id := c.Param("user_id")

	rows, err := conn.QueryContext(c.Request.Context(), "select * from tweets_stats where user_id = ?;", user_id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Tweets found"})
		return
	} else if err != nil {
		log.Panic("ERROR_GETTING_TWEETS_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	defer rows.Close()
	var tweets []models.GetTweetsResponse
	for rows.Next() {
		var tweet models.GetTweetsResponse
		if err := rows.Scan(&tweet.Tweet_id, &tweet.Body, &tweet.User_id, &tweet.Created_At, &tweet.Tweet_total_likes); err != nil {
			log.Printf("ERROR_SCANNING_TWEET_ROW: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning tweet data"})
			return
		}
		tweets = append(tweets, tweet)
	}

	// 🔥 Store into Redis cache
	// 🧹 Marshal struct to JSON before caching
	tweetsJSON, err := json.Marshal(tweets)
	if err != nil {
		log.Println("Failed to marshal user:", err)
	} else {
		cacheKey := c.Request.URL.Path
		err := rdb.Set(c, cacheKey, tweetsJSON, 5*time.Minute).Err()
		if err != nil {
			log.Println("Redis cache SET error:", err)
		}
	}

	c.JSON(http.StatusOK, tweets)
}

func CreateTweet(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()
	var req models.CreateTweetRequest
	user_id := c.Param("user_id")
	// Bind incoming JSON to the struct
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err = conn.ExecContext(c.Request.Context(), "insert into tweets (body, user_id) values (?, ?);", req.Body, user_id)
	if err != nil {
		log.Panic("ERROR_CREATING_TWEET: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	log.Println("tweet created by user_id:", user_id)
	c.Status(http.StatusCreated)
}

func UpdateTweet(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()
	var req models.UpdateTweetRequest
	user_id := c.Param("user_id")
	tweet_id := c.Param("tweet_id")
	// Bind incoming JSON to the struct
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err = conn.ExecContext(c.Request.Context(), "update tweets SET body = ? where tweet_id = ? and user_id = ?;", req.Body, tweet_id, user_id)
	if err != nil {
		log.Panic("ERROR_UPDATING_TWEET: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	log.Println("updated tweet tweet_id:", tweet_id)
	c.Status(http.StatusNoContent)
}

func DeleteTweet(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()
	user_id := c.Param("user_id")
	tweet_id := c.Param("tweet_id")
	_, err = conn.ExecContext(c.Request.Context(), "delete from tweets where user_id = ? and tweet_id = ?;", user_id, tweet_id)
	if err != nil {
		log.Panic("ERROR_DELETING_TWEET_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Println("tweet deleted user_id:", user_id)
	c.Status(http.StatusNoContent)
}

func LikeTweet(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()
	user_id := c.Param("user_id")
	tweet_id := c.Param("tweet_id")
	_, err = conn.ExecContext(c.Request.Context(), "insert into likes_tweets (tweet_id, user_id) values (?,?);", tweet_id, user_id)
	if err != nil {
		log.Panic("ERROR_LIKING_TWEET_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("tweet (%s) liked by user_id: (%s)", tweet_id, user_id)
	c.Status(http.StatusNoContent)
}

func UnLikeTweet(c *gin.Context, db *sql.DB) {
	// get db connection and release it when the transaction is complete
	conn, err := dbHelper.GetDBConn(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}
	defer conn.Close()
	user_id := c.Param("user_id")
	tweet_id := c.Param("tweet_id")
	_, err = conn.ExecContext(c.Request.Context(), "delete from likes_tweets where user_id = ? and tweet_id = ?;", user_id, tweet_id)
	if err != nil {
		log.Panic("ERROR_DELETING_LIKE_DATA: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	log.Printf("tweet (%s) unliked by user_id: (%s)", tweet_id, user_id)
	c.Status(http.StatusNoContent)
}