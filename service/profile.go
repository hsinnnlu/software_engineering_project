package service

import (
	"github.com/hsinnnlu/software_engineering_project/db"
	"github.com/hsinnnlu/software_engineering_project/models"
)

func GetAnnouncementList() ([]models.Announce, error) {
	announces := []models.Announce{}
	query := "SELECT announce_id, announce_title, announce_content, announce_date FROM Announcements"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		announce := models.Announce{}
		err := rows.Scan(&announce.Announce_id, &announce.Announce_title, &announce.Announce_content, &announce.Announce_date)
		if err != nil {
			return nil, err
		}
		announces = append(announces, announce)
	}
	return announces, nil
}
