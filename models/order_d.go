package models

import "time"

type OrderD struct {
	OrderDID    int     `json:"order_d_id" gorm:"primary_key;auto_increment:true"`
	BarberID    int     `json:"barber_id" gorm:"type:integer"`
	OrderID     int     `json:"order_id" gorm:"primary_key;type:integer"`
	PaketID     int     `json:"paket_id" gorm:"type:integer;not null"`
	PaketName   string  `json:"paket_name" gorm:"type:varchar(60)"`
	DurasiStart int     `json:"durasi_start" gorm:"type:integer"`
	DurasiEnd   int     `json:"durasi_end" gorm:"type:integer"`
	Price       float32 `json:"price" gorm:"type:numeric(20,4)"`
	Model
}

type OrderDPost struct {
	PaketID     int     `json:"paket_id" valid:"Required"`
	PaketName   string  `json:"paket_name" valid:"Required"`
	DurasiStart int     `json:"durasi_start" valid:"Required"`
	DurasiEnd   int     `json:"durasi_end" valid:"Required"`
	Price       float32 `json:"price" valid:"Required"`
}
type OrderHGet struct {
	BarberID    int       `json:"barber_id"`
	BarberCd    string    `json:"barber_cd"`
	OrderDate   time.Time `json:"order_date"`
	BarberName  string    `json:"barber_name"`
	CapsterID   int       `json:"capster_id"`
	CapsterName string    `json:"capster_name"`
	FileID      int       `json:"file_id" `
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	FromApps    bool      `json:"from_apps"`
	Status      string    `json:"status"`
	Price       float32   `json:"price"`
}

type OrderDGet struct {
	BarberName  string  `json:"barber_name"`
	CapsterID   int     `json:"capster_id"`
	CapsterName string  `json:"capster_name"`
	FileID      int     `json:"file_id" `
	FileName    string  `json:"file_name"`
	FilePath    string  `json:"file_path"`
	PaketID     int     `json:"paket_id"`
	PaketName   string  `json:"paket_name"`
	Price       float32 `json:"price"`
	DurasiStart int     `json:"durasi_start"`
	DurasiEnd   int     `json:"durasi_end"`
}
