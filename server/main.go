package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/treavorgagne/twitter-clone/server/config"
	"github.com/treavorgagne/twitter-clone/server/routes"
)

func HealthCheck(c *gin.Context) {
    c.Status(http.StatusNoContent);
}

func main() {
	db := config.ConfigDB()
    defer db.Close()
    log.Println("db connection pool opened")

    router := gin.Default()
    router.SetTrustedProxies(nil)
    log.Println("started gin router")

    rdb := config.CacheConn()
    defer rdb.Close()

    // healh route
    router.GET("/health", HealthCheck);

    // check cache
    router.Use(config.CacheMiddleware(rdb))

    // set new db connection
    router.Use(config.GetDBConn(db))

    // user routes
    router.POST("/users", routes.CreateUser)
    router.GET("/users/:user_id", routes.GetUser)
    router.PUT("/users/:user_id", routes.UpdateUser);
    router.DELETE("/users/:user_id", routes.DeleteUser);

    // follows routes
    router.POST("/users/:user_id/follows/:follow_id", routes.FollowUser);
    router.DELETE("/users/:user_id/follows/:follow_id", routes.UnFollowUser);

    // tweets routes
    router.GET("/tweets/:tweet_id", routes.GetTweet);
    router.GET("/users/:user_id/tweets", routes.GetTweets);
    router.POST("/users/:user_id/tweets", routes.CreateTweet);
    router.PUT("/users/:user_id/tweets/:tweet_id", routes.UpdateTweet);
    router.DELETE("/users/:user_id/tweets/:tweet_id", routes.DeleteTweet);
    router.POST("/users/:user_id/tweets/:tweet_id/likes", routes.LikeTweet);
    router.DELETE("/users/:user_id/tweets/:tweet_id/unlikes", routes.UnLikeTweet);

    // comments routes
    router.GET("/tweets/:tweet_id/comments", routes.GetCommentsByTweetId);
    router.POST("/users/:user_id/tweets/:tweet_id/comment", routes.CreateComment);
    router.GET("comment/:comment_id", routes.GetComment);
    router.PUT("comment/:comment_id", routes.UpdateComment);
    router.DELETE("/comment/:comment_id", routes.DeleteComment);
    router.POST("/users/:user_id/comment/:comment_id/like", routes.LikeComment);
    router.DELETE("/users/:user_id/comment/:comment_id/unlike", routes.UnLikeComment);

    router.Run("0.0.0.0:" + os.Getenv("SERVERPORT"))
}