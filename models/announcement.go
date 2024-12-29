// 通知的頁面
package models

type Announce struct {
	Announce_id		 int    `json:"id"`
	Announce_title   string `json:"title"`
	Announce_content string `json:"content"`
	Announce_date    string `json:"date"`
}
