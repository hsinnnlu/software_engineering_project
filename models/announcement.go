package models

type Announce struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    string `json:"time"`
}
