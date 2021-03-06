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
	OrderID       int       `json:"order_id"`
	OrderNo       string    `json:"order_no"`
	Status        string    `json:"status"`
	OrderDate     time.Time `json:"order_date"`
	BarberID      int       `json:"barber_id"`
	BarberName    string    `json:"barber_name"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Distance      float32   `json:"distance"`
	BarberRating  float32   `json:"barber_rating"`
	CapsterID     int       `json:"capster_id"`
	CapsterName   string    `json:"capster_name"`
	CapsterRating float32   `json:"capster_rating"`
	SaFileOutput
	PaketID    int     `json:"paket_id"`
	TotalPrice float32 `json:"total_price"`
}

type OrderDataBy struct {
	OrderDGet
	DataDetail         []*OrderD       `json:"data_detail"`
	DataFeedbackRating *FeedbackRating `json:"data_feedback_rating"`
}

// DataDetail         []*OrderD       `json:"data_detail"`
// 	DataFeedbackRating *FeedbackRating `json:"data_feedback_rating"`
