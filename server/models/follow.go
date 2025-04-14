package models

type FollowRequest struct {
	User_id 		int 	`form:"user_id" json:"user_id"`
	Follow_id 		int 	`form:"follow_id" json:"follow_id"`
}

type UnFollowRequest struct {
	User_id 		int 	`form:"user_id" json:"user_id"`
	Follow_id 		int 	`form:"follow_id" json:"follow_id"`
}