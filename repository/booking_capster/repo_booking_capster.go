package repobookingcapster

import (
	"fmt"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"

	ibookingcapster "nuryanto2121/cukur_in_user/interface/booking_capster"

	"github.com/jinzhu/gorm"
)

type repoBookingCapster struct {
	Conn *gorm.DB
}

func NewRepoBookingCapster(Conn *gorm.DB) ibookingcapster.Repository {
	return &repoBookingCapster{Conn}
}
func (db *repoBookingCapster) GetDataBy(param models.AddBookingCapster) (result *[]models.BookingCapster, err error) {
	var (
		// pageNum  = 0
		// pageSize = setting.FileConfigSetting.App.PageSize
		sWhere = ""
		logger = logging.Logger{}
		query  *gorm.DB
	)
	if param.BookingDate.IsZero() {
		query = db.Conn.Where(`barber_id = ? and capster_id = ?`, param.BarberID, param.CapsterID).Order(`booking_date`).Find(&result)
	} else {
		query = db.Conn.Where(`barber_id = ? and capster_id = ? AND booking_date = ?`, param).Order(`booking_date`).Find(&result)
	}
	query = db.Conn.Where(sWhere).Order(`booking_date`).Find(&result)

	// end where

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (db *repoBookingCapster) Create(data *models.BookingCapster) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Create(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
