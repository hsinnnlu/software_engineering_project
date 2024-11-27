// 2024.11.26 新增、查詢講座資訊 finished
package models

type Lecture struct {
	Lecture_id      	int    `json:"id"`
	Lecture_name    	string `json:"name"`
	Lecture_speaker 	string `json:"speaker"`
	Lecture_timestamp	string `json:"timestamp"`
	Lecture_manager     string `json:"manager"`
	Qty_participater 	int	   `json:"qty_participater"`
	Lecture_location    string `json:"location"`
	Status          	string `json:"status"`
}