package models

import "time"

type Capster struct {
	Name          string    `json:"name" valid:"Required"`
	Telp          string    `json:"telp,omitempty"`
	Email         string    `json:"email"  valid:"Required"`
	UserType      string    `json:"user_type" valid:"Required"`
	IsActive      bool      `json:"is_active" valid:"Required"`
	JoinDate      time.Time `json:"join_date" valid:"Required"`
	FileID        int       `json:"file_id,omitempty"`
	TopCollection []Collections
}

type CapsterUpdate struct {
	Name     string    `json:"name" valid:"Required"`
	Telp     string    `json:"telp,omitempty"`
	Email    string    `json:"email"  valid:"Required"`
	UserType string    `json:"user_type" valid:"Required"`
	IsActive bool      `json:"is_active" valid:"Required"`
	JoinDate time.Time `json:"join_date" valid:"Required"`
	FileID   int       `json:"file_id,omitempty"`
}

type CapsterCollection struct {
	CapsterID int       `json:"capster_id" gorm:"primary_key;type:integer"`
	FileID    int       `json:"file_id" gorm:"primary_key;type:integer"`
	UserInput string    `json:"user_input" gorm:"type:varchar(20)"`
	UserEdit  string    `json:"user_edit" gorm:"type:varchar(20)"`
	TimeInput time.Time `json:"time_input" gorm:"type:timestamp(0) without time zone;default:now()"`
	TimeEdit  time.Time `json:"time_Edit" gorm:"type:timestamp(0) without time zone;default:now()"`
}

// type CapsterCollectionPath struct {
// 	CapsterCollection
// 	FileName string `json:"file_name"`
// 	FilePath string `json:"file_path"`
// }

type Collections struct {
	FileID int `json:"file_id"`
}

type CapsterList struct {
	CapsterID      int       `json:"capster_id" valid:"Required"`
	CapsterName    string    `json:"capster_name,omitempty"`
	CapsterRating  float32   `json:"capster_rating"`
	IsBusy         bool      `json:"is_busy"`
	IsActive       bool      `json:"is_active" valid:"Required" gorm:"type:boolean"`
	JoinDate       time.Time `json:"join_date"`
	BarberID       int       `json:"barber_id"`
	BarberName     string    `json:"barber_name"`
	BarberRating   float32   `json:"barber_rating"`
	Distance       float32   `json:"distance"`
	IsBarberOpen   bool      `json:"is_barber_open" gorm:"type:boolean"`
	IsBarberActive bool      `json:"is_barber_active" gorm:"type:boolean"`
	OperationStart time.Time `json:"operation_start"`
	OperationEnd   time.Time `json:"operation_end"`
	LengthOfWork   string    `json:"length_of_work"`
	SaFileOutput
}
