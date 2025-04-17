package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	dbHelper "github.com/treavorgagne/twitter-clone/server/db"
	"github.com/treavorgagne/twitter-clone/server/redis"
	"github.com/treavorgagne/twitter-clone/server/routes"
)

func HealthCheck(c *gin.Context) {
    c.Status(http.StatusNoContent);
}

func main() {
	db := dbHelper.ConfigDB()
    defer db.Close()
    log.Println("db connection pool opened")

    router := gin.Default()
    router.SetTrustedProxies(nil)
    log.Println("started gin router")

    rdb := redis.CacheConn()
    defer rdb.Close()

    // healh route
    router.GET("/health", HealthCheck);

    // check cache
    router.Use(redis.CacheMiddleware(rdb))

    // user routes
    router.POST("/users",               func(c *gin.Context) { routes.CreateUser(c, db)})
    router.GET("/users/:user_id",       func(c *gin.Context) { routes.GetUser(c, db, rdb) })
    router.PUT("/users/:user_id",       func(c *gin.Context) { routes.UpdateUser(c, db) });
    router.DELETE("/users/:user_id",    func(c *gin.Context) { routes.DeleteUser(c, db) });

    // follows routes
    router.POST("/users/:user_id/follows/:follow_id",       func(c *gin.Context) { routes.FollowUser(c, db) });
    router.DELETE("/users/:user_id/follows/:follow_id",     func(c *gin.Context) { routes.UnFollowUser(c, db) });

    // tweets routes
    router.GET("/tweets/:tweet_id",                             func(c *gin.Context) { routes.GetTweet(c, db, rdb) });
    router.GET("/users/:user_id/tweets",                        func(c *gin.Context) { routes.GetTweets(c, db, rdb) });
    router.POST("/users/:user_id/tweets",                       func(c *gin.Context) { routes.CreateTweet(c, db) });
    router.PUT("/users/:user_id/tweets/:tweet_id",              func(c *gin.Context) { routes.UpdateTweet(c, db) });
    router.DELETE("/users/:user_id/tweets/:tweet_id",           func(c *gin.Context) { routes.DeleteTweet(c, db) });
    router.POST("/users/:user_id/tweets/:tweet_id/likes",       func(c *gin.Context) { routes.LikeTweet(c, db) });
    router.DELETE("/users/:user_id/tweets/:tweet_id/unlikes",   func(c *gin.Context) { routes.UnLikeTweet(c, db) });

    // // comments routes
    router.GET("/tweets/:tweet_id/comments",                      func(c *gin.Context) { routes.GetCommentsByTweetId(c, db, rdb) });
    router.POST("/users/:user_id/tweets/:tweet_id/comment",       func(c *gin.Context) { routes.CreateComment(c, db) });
    router.GET("comment/:comment_id",                             func(c *gin.Context) { routes.GetComment(c, db, rdb) });
    router.PUT("comment/:comment_id",                             func(c *gin.Context) { routes.UpdateComment(c, db) });
    router.DELETE("/comment/:comment_id",                         func(c *gin.Context) { routes.DeleteComment(c, db) });
    router.POST("/users/:user_id/comment/:comment_id/like",       func(c *gin.Context) { routes.LikeComment(c, db) });
    router.DELETE("/users/:user_id/comment/:comment_id/unlike",   func(c *gin.Context) { routes.UnLikeComment(c, db) });

    router.Run("localhost:8080")
}