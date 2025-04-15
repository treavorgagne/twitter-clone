package models

type GetTweetsResponse struct {
	Tweet_id       		int    	`json:"tweet_id"`
	Body 				string 	`json:"body"`
	User_id       		int    	`json:"user_id"`
	Created_At			string  `json:"created_at"`
	Tweet_total_likes 	int 	`json:"tweet_total_likes"`
}

type CreateTweetRequest struct {
	Body string `form:"body" json:"body"`
}

type UpdateTweetRequest struct {
	Body string `form:"body" json:"body"`
}