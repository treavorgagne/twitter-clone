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

    // // follows routes
    // router.POST("/users/:user_id/follows/:following_id", routes.FollowUser);
    // router.DELETE("/users/:user_id/follows/:following_id", routes.UnFollowUser);

    // // tweets routes
    // router.GET("/users/:user_id/tweets", routes.GetTweet); 
    // router.POST("/users/:user_id/tweets", routes.CreateTweet);
    // router.DELETE("/users/:user_id/tweets/:tweet_id", routes.DeleteTweet);
    // router.POST("/users/:user_id/tweets/:tweet_id/like", routes.LikeTweet);
    // router.DELETE("/users/:user_id/tweets/:tweet_id/unlike", routes.UnLikeTweet);

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