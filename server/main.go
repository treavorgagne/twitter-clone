package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/treavorgagne/twitter-clone/server/routes"
)

func HealthCheck(c *gin.Context) {
    c.Status(http.StatusNoContent);
}

func main() {
    router := gin.Default()
    router.SetTrustedProxies(nil)

    // healh route
    router.GET("/health", HealthCheck);

    // user routes
    router.POST("/users", routes.CreateUser);
    router.GET("/users/:user_id", routes.GetUser);
    router.PUT("/users/:user_id", routes.UpdateUser);
    router.DELETE("/users/:user_id", routes.DeleteUser);

    // follows routes
    router.POST("/users/:user_id/follows/:following_id", routes.FollowUser);
    router.DELETE("/users/:user_id/follows/:following_id", routes.UnFollowUser);

    // tweets routes
    router.POST("/users/:user_id/tweets", routes.CreateTweet);
    router.DELETE("/users/:user_id/tweets/:tweet_id", routes.DeleteTweet);
    router.POST("/users/:user_id/tweets/:tweet_id/like", routes.LikeTweet);
    router.DELETE("/users/:user_id/tweets/:tweet_id/unlike", routes.UnLikeTweet);

    // comments routes
    router.POST("/users/:user_id/tweets/:tweet_id/comment", routes.CreateComment);
    router.DELETE("/users/:user_id/tweets/:tweet_id/comment/:comment_id", routes.DeleteComment);
    router.POST("/users/:user_id/tweets/:tweet_id/comment/:comment_id/like", routes.LikeComment);
    router.DELETE("/users/:user_id/tweets/:tweet_id/comment/:comment_id/unlike", routes.UnLikeComment);

    router.Run("localhost:8080")
}