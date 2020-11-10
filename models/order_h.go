package models

import "time"

type OrderH struct {
	OrderID      int       `json:"order_id" gorm:"primary_key;auto_increment:true"`
	OrderNo      string    `json:"order_no" gorm:"type:varchar(20)"`
	BarberID     int       `json:"barber_id" gorm:"type:integer"`
	CapsterID    int       `json:"capster_id" gorm:"type:integer"`
	OrderDate    time.Time `json:"order_date" gorm:"type:timestamp(0) without time zone"`
	UserID       int       `json:"user_id" gorm:"type:integer"`
	CustomerName string    `json:"customer_name" gorm:"type:varchar(60);not null"`
	Telp         string    `json:"telp" gorm:"type:varchar(20)"`
	Status       string    `json:"status" gorm:"type:varchar(1)"`
	FromApps     bool      `json:"from_apps" gorm:"type:boolean"`
	Model
}

type OrderPost struct {
	BarberID     int          `json:"barber_id" valid:"Required"`
	CapsterID    int          `json:"capster_id,omitempty"`
	OrderDate    time.Time    `json:"order_date" valid:"Required"`
	UserID       int          `json:"user_id,omitempty"`
	CustomerName string       `json:"customer_name" valid:"Required"`
	Telp         string       `json:"telp,omitempty"`
	Pakets       []OrderDPost `json:"paket_ids"`
}

type OrderList struct {
	OwnerID     int       `json:"owner_id"`
	BarberID    int       `json:"barber_id" valid:"Required"`
	BarberName  string    `json:"barber_name"`
	OrderID     int       `json:"order_id"`
	Status      string    `json:"status" `
	FromApps    bool      `json:"from_apps"`
	CapsterID   int       `json:"capster_id,omitempty"`
	CapsterName string    `json:"capster_name"`
	OrderDate   time.Time `json:"order_date" valid:"Required"`
	FileID      int       `json:"file_id" `
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	Price       float32   `json:"price" `
	Weeks       int       `json:"weeks"`
	Months      int       `json:"months"`
	Years       int       `json:"years"`
}
