package repobookingcapster

import (
	"fmt"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"

	ibookingcapster "nuryanto2121/cukur_in_user/interface/booking_capster"

	"gorm.io/gorm"
)

type repoBookingCapster struct {
	Conn *gorm.DB
}

func NewRepoBookingCapster(Conn *gorm.DB) ibookingcapster.Repository {
	return &repoBookingCapster{Conn}
}
func (db *repoBookingCapster) GetDataBy(param models.AddBookingCapster) (result *[]models.OrderH, err error) {
	var (
		// pageNum  = 0
		// pageSize = setting.FileConfigSetting.App.PageSize
		sWhere = ""
		Rest   = &[]models.OrderH{}
		logger = logging.Logger{}
		query  *gorm.DB
	)

	sWhere = fmt.Sprintf(`barber_id = ? and capster_id = ? AND order_date::date = '%v'::date AND status = 'N'`, param.BookingDate.Format("2006-01-02"))
	query = db.Conn.Where(sWhere, param.BarberID, param.CapsterID).Order(`order_date`).Find(&Rest)

	// query = db.Conn.Where(sWhere).Order(`booking_date`).Find(&result)

	// end where

	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return Rest, nil
}

func (db *repoBookingCapster) Count(param models.AddBookingCapster, UserID int) (result int, err error) {
	var (
		sWhere  = ""
		logger  = logging.Logger{}
		query   *gorm.DB
		_result int64 = 0
	)
	result = 0

	// sWhere = fmt.Sprintf(`barber_id = ? and capster_id = ? AND order_date = ? AND user_id = ? and status = 'N'`)
	sWhere = fmt.Sprintf(`barber_id = ? and capster_id = ? AND user_id = ? AND order_date::date = '%v'::date AND status = 'N'`, param.BookingDate.Format("2006-01-02"))
	query = db.Conn.Table(`order_h`).Where(sWhere, param.BarberID, param.CapsterID, UserID).Count(&_result)

	// query = db.Conn.Where(sWhere).Order(`booking_date`).Find(&result)

	// end where

	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error

	if err != nil {
		return 0, err
	}
	return int(_result), nil
}

func (db *repoBookingCapster) Create(data *models.BookingCapster) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Create(data)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
