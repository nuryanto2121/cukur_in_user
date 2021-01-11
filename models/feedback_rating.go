package models

type FeedbackRating struct {
	ID     int `json:"id" gorm:"primary_key;auto_increment:true"`
	UserID int `json:"user_id" gorm:"type:integer;not null"`
	AddFeedbackRating
	Model
}

type AddFeedbackRating struct {
	BarberID      int     `json:"barber_id" valid:"Required" gorm:"type:integer;not null"`
	CapsterID     int     `json:"capster_id" valid:"Required" gorm:"type:integer;not null"`
	OrderID       int     `json:"order_id" valid:"Required" gorm:"type:integer;not null"`
	BarberRating  float32 `json:"barber_rating" gorm:"type:numeric(5,2);default:5"`
	CapsterRating float32 `json:"capster_rating" gorm:"type:numeric(5,2);default:5"`
	Comment       string  `json:"comment" gorm:"type:varchar(1000)"`
}

type OutFeedbackRating struct {
	BarberRating  float32 `json:"barber_rating" gorm:"type:numeric(5,2);default:5"`
	CapsterRating float32 `json:"capster_rating" gorm:"type:numeric(5,2);default:5"`
	Comment       string  `json:"comment" gorm:"type:varchar(1000)"`
}
