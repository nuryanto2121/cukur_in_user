package models

type SsUserFcm struct {
	SsUserFcmID int    `json:"ss_user_fcm_id" gorm:"PRIMARY_KEY"`
	UserID      int    `json:"user_id" gorm:"type:integer;not null"`
	FcmToken    string `json:"fcm_token" gorm:"type:varchar(150)"`
	Token       string `json:"token" gorm:"type:varchar(150)"`
	Model
}
