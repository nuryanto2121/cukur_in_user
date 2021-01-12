package repobookingcapster

import (
	"fmt"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"time"

	ibookingcapster "nuryanto2121/cukur_in_user/interface/booking_capster"

	"gorm.io/gorm"
)

type repoBookingCapster struct {
	Conn *gorm.DB
}

func NewRepoBookingCapster(Conn *gorm.DB) ibookingcapster.Repository {
	return &repoBookingCapster{Conn}
}
func (db *repoBookingCapster) GetDataBy(param models.AddBookingCapster) (result interface{}, err error) {
	type Results struct {
		ScheduleTime time.Time `json:"schedule_time"`
	}

	var (
		// pageNum  = 0
		// pageSize = setting.FileConfigSetting.App.PageSize
		op = []*Results{}
		// Rest   = &[]models.OrderH{}
		logger = logging.Logger{}
		query  *gorm.DB
	)

	sSql := fmt.Sprintf(`
		SELECT h.schedule_time as schedule_time
		FROM   (select
					generate_series (b2.operation_start::timestamp 
							, b2.operation_end::timestamp 
							, interval '30m') as schedule_time,
					b2.operation_start ,
					b2.operation_end ,
					b2.barber_id 
				from barber b2 
				where b2.barber_id = %d
				)h 
		WHERE  EXTRACT(ISODOW FROM h.schedule_time) < 6
		and h.schedule_time::time not in (
			select order_date::time from order_h  
			where barber_id = %d and capster_id = %d AND order_date::date = '%v'  AND status = 'N'
		)
		and 	h.schedule_time::time >= now()::time
		AND   now()::time >= h.operation_start::time 
		AND   now()::time <= h.operation_end::time ;
	`, param.BarberID, param.BarberID, param.CapsterID, param.BookingDate.Format("2006-02-01"))
	fmt.Println(sSql)
	// sWhere = fmt.Sprintf(`barber_id = ? and capster_id = ? AND order_date::date = '%v'::date AND status = 'N'`, param.BookingDate.Format("2006-01-02"))
	// query = db.Conn.Where(sWhere, param.BarberID, param.CapsterID).Order(`order_date`).Find(&Rest)

	query = db.Conn.Raw(sSql).Find(&op)

	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string

	err = query.Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return op, nil
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
