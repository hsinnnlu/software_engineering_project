// 2024.11.26 學生的聽講狀況的struct 
package models

type User_Lecture struct {
	User          string	`json: user`
	Lecture       int		`json: lecture`
	Sign_in_time  string	`json: sign_in_time`
	Sign_out_time string	`json: sign_out_time`
}
// time 轉換