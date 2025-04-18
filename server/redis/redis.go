package redis

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// CacheConn initializes and returns a Redis client
func CacheConn() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:4444",
		DB:       0, // optional: default DB
	})
	return rdb
}

func CacheMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		cacheKey := c.Request.URL.Path
		val, err := rdb.Get(c, cacheKey).Result()
		if err == redis.Nil { // Cache miss
			log.Println("Redis miss. Storing key: ", cacheKey)
			c.Set("rdb", rdb);
			c.Next()
		} else if err != nil { // Redis error
			log.Println("Redis error:", err)
			c.Next()
		} else { // Cache hit
			log.Println("Redis hit for key:", cacheKey)
			c.Data(200, "application/json", []byte(val))
			c.Abort()
		}
	}
}