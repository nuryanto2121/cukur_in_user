package models

import "time"

type Advertis struct {
	AdvertisID int `json:"advertis_id" gorm:"PRIMARY_KEY"`
	AddAdvertis
	Model
}

type AddAdvertis struct {
	Title          string    `json:"title" gorm:"type:varchar(100)"`
	Descs          string    `json:"descs" gorm:"type:varchar(255)"`
	AdvertisStatus string    `json:"advertis_status" gorm:"type:varchar(2)"`
	SlideDuration  int       `json:"slide_duration" gorm:"type:integer"`
	StartDate      time.Time `json:"start_date" gorm:"type:timestamp(0) without time zone"`
	EndDate        time.Time `json:"end_date" gorm:"type:timestamp(0) without time zone"`
	FileID         int       `json:"file_id" gorm:"primary_key;auto_increment:true"`
}

type ListAdvertis struct {
	AdvertisID int `json:"advertis_id"`
	AddAdvertis
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
}
