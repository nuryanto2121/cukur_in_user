package models

import "time"

type Advertis struct {
	AdvertisID     int       `json:"advertis_id" gorm:"PRIMARY_KEY"`
	Title          string    `json:"title" gorm:"varchar(100)"`
	Descs          string    `json:"descs" gorm:"varchar(255)"`
	AdvertisStatus string    `json:"advertis_status" gorm:"type:varchar(2)"`
	SlideDuration  int       `json:"slide_duration" gorm:"type:integer"`
	StartDate      time.Time `json:"start_date" gorm:"type:timestamp(0)"`
	EndDate        time.Time `json:"end_date" gorm:"type:timestamp(0)"`
	FileName       string    `json:"file_name"`
	FilePath       string    `json:"file_path"`
	FileType       string    `json:"file_type"`
	Model
}
