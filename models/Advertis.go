package models

import "time"

type Advertise struct {
	AdvertiseID int `json:"advertise_id" gorm:"PRIMARY_KEY"`
	AddAdvertise
	Model
}

type AddAdvertise struct {
	Title           string    `json:"title" gorm:"type:varchar(100)"`
	Descs           string    `json:"descs" gorm:"type:varchar(255)"`
	AdvertiseStatus string    `json:"advertise_status" gorm:"type:varchar(2)"`
	SlideDuration   int       `json:"slide_duration" gorm:"type:integer"`
	StartDate       time.Time `json:"start_date" gorm:"type:timestamp(0) without time zone"`
	EndDate         time.Time `json:"end_date" gorm:"type:timestamp(0) without time zone"`
	FileID          int       `json:"file_id" gorm:"primary_key;auto_increment:true"`
}

type ListAdvertise struct {
	AdvertiseID int `json:"advertise_id"`
	AddAdvertise
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileType string `json:"file_type"`
}
