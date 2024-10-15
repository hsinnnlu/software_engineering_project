package models

type Lecture struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Speaker         string `json:"speaker"`
	Timestamp       string `json:"timestamp"`
	Manager         string `json:"manager"`
	Qty_participate string `json:"qty_participate"`
	Location        string `json:"location"`
	Status          string `json:"status"`
}
