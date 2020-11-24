package models

import "time"

type BookingCapster struct {
	BookingCapsterID int `json:"booking_capster_id" gorm:"primary_key;auto_increment:true"`
	AddBookingCapster
	Model
}

type AddBookingCapster struct {
	BarberID    int       `json:"barber_id" gorm:"type:integer"`
	CapsterID   int       `json:"capster_id" gorm:"type:integer"`
	BookingDate time.Time `json:"booking_date" gorm:"type:timestamp(0) without time zone"`
}
