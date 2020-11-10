package models

type FeedbackRating struct {
	BarberID      int     `json:"barber_id" gorm:"primary_key;auto_increment:false"`
	CapsterID     int     `json:"capster_id" gorm:"primary_key;auto_increment:false"`
	UserID        int     `json:"user_id" gorm:"primary_key;auto_increment:false"`
	BarberRating  float32 `json:"barber_rating" gorm:"type:numeric(5,2)"`
	CapsterRating float32 `json:"capster_rating" gorm:"type:numeric(5,2)"`
	Model
}
