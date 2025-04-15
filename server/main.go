package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/treavorgagne/twitter-clone/server/db"
	"github.com/treavorgagne/twitter-clone/server/routes"
)

func main() {
	var db, err = db.ConnectToDB()
	if err != nil {
		log.Fatal("ERROR_DB_CONNECTION", err)
	}

    router := gin.Default()
    router.SetTrustedProxies(nil)

    // healh route
    router.GET("/health", HealthCheck);

    // user routes
    router.POST("/users", routes.CreateUser(db));
    router.GET("/users/:user_id", routes.GetUser(db));
    router.PUT("/users/:user_id", routes.UpdateUser(db));
    router.DELETE("/users/:user_id", routes.DeleteUser(db));

    // follows routes
    router.POST("/users/:user_id/follows/:follow_id", routes.FollowUser(db));
    router.DELETE("/users/:user_id/follows/:follow_id", routes.UnFollowUser(db));

    // // tweets routes
    router.GET("/tweets/:tweet_id", routes.GetTweet(db));
    router.GET("/users/:user_id/tweets", routes.GetTweets(db));
    router.POST("/users/:user_id/tweets", routes.CreateTweet(db));
    router.PUT("/users/:user_id/tweets/:tweet_id", routes.UpdateTweet(db));
    router.DELETE("/users/:user_id/tweets/:tweet_id", routes.DeleteTweet(db));
    router.POST("/users/:user_id/tweets/:tweet_id/likes", routes.LikeTweet(db));
    router.DELETE("/users/:user_id/tweets/:tweet_id/unlikes", routes.UnLikeTweet(db));

    // // comments routes
    // router.POST("/users/:user_id/tweets/:tweet_id/comment", routes.CreateComment);
    // router.DELETE("/users/:user_id/tweets/:tweet_id/comment/:comment_id", routes.DeleteComment);
    // router.POST("/users/:user_id/tweets/:tweet_id/comment/:comment_id/like", routes.LikeComment);
    // router.DELETE("/users/:user_id/tweets/:tweet_id/comment/:comment_id/unlike", routes.UnLikeComment);

    router.Run("localhost:8080")

}

func HealthCheck(c *gin.Context) {
    c.Status(http.StatusNoContent);
}