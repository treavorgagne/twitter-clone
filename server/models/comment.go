package models

type GetCommentsResponse struct {
	Comment_id       		int    	`json:"comment_id"`
	Tweet_id       			int    	`json:"tweet_id"`
	Body 					string 	`json:"body"`
	User_id       			int    	`json:"user_id"`
	Created_At				string  `json:"created_at"`
	Comment_total_likes 	int 	`json:"comment_total_likes"`
}

type CreateCommentRequest struct {
	Body string `form:"body" json:"body"`
}

type UpdateCommentRequest struct {
	Body string `form:"body" json:"body"`
}