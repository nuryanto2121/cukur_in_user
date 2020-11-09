package version

import "github.com/jinzhu/gorm"

type SsVersion struct {
	VersionID int    `json:"version_id" gorm:"PRIMARY_KEY"`
	OS        string `json:"os" gorm:"type:varchar(20)"`
	Version   int    `json:"version" gorm:"type:integer"`
	Apps      string `json:"apps" gorm:"type:varchar(20)"`
}

func (V *SsVersion) GetVersion(Conn *gorm.DB) (result SsVersion, err error) {
	err = Conn.Where("os = ? AND apps = 'user' ", V.OS).First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}
