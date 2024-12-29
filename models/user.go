// 2024.11.26 新增帳號的struct finished
package models

type User struct {
	User_id       	 string  `json:"user_id"`
	Password_hash 	 string  `json: password_hash`
	User_name     	 string  `json: user_name`
	User_email    	 string  `json: user_email`
	User_permission	 int  	 `json: permission`
	Qty_lecture		 int	 `json: qty_lecture`
	User_professor   string  `json: user_professor`
}
