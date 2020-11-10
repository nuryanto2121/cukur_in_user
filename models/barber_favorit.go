package models

type BarberFavorit struct {
	AddBarberFavorit
	Model
}

type BarberFavoritList struct {
	BarbersList
	BarberRating float32 `json:"barber_rating"`
}

type AddBarberFavorit struct {
	BarberID int `json:"barber_id" gorm:"primary_key;auto_increment:false"`
	UserID   int `json:"user_id" gorm:"PRIMARY_KEY;auto_increment:false"`
}
