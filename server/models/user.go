package models

type CreateUserRequest struct {
	Username string `form:"username" json:"username"`
}

type UpdateUserRequest struct {
	Username string `form:"username" json:"username"`
}

type GetUserResponse struct {
	User_id       		int    	`json:"user_id"`
	Username 			string 	`json:"username"`
	Created_At			string  `json:"created_at"`
	Total_following    	int 	`json:"total_following"`
	Total_followers	   	int 	`json:"total_followers"` 	
}


