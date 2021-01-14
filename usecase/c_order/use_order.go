package useorder

import (
	"context"
	"errors"
	"fmt"
	"math"
	ibarber "nuryanto2121/cukur_in_user/interface/barber"
	ibookingcapster "nuryanto2121/cukur_in_user/interface/booking_capster"
	iorderd "nuryanto2121/cukur_in_user/interface/c_order_d"
	iorderh "nuryanto2121/cukur_in_user/interface/c_order_h"
	ifeedbackrating "nuryanto2121/cukur_in_user/interface/feedback_rating"
	inotification "nuryanto2121/cukur_in_user/interface/notification"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"nuryanto2121/cukur_in_user/redisdb"
	repofunction "nuryanto2121/cukur_in_user/repository/function"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

type useOrder struct {
	repoOrderH         iorderh.Repository
	repoOrderD         iorderd.Repository
	repoBarber         ibarber.Repository
	repoBookingCapster ibookingcapster.Repository
	usenotification    inotification.Usecase
	repoNotif          inotification.Repository
	repoFeedbackRating ifeedbackrating.Repository
	contextTimeOut     time.Duration
}

func NewUserMOrder(a iorderh.Repository, b iorderd.Repository, c ibarber.Repository,
	d ibookingcapster.Repository, e inotification.Usecase,
	f ifeedbackrating.Repository, g inotification.Repository,
	timeout time.Duration) iorderh.Usecase {
	return &useOrder{
		repoOrderH:         a,
		repoOrderD:         b,
		repoBarber:         c,
		repoBookingCapster: d,
		usenotification:    e,
		repoFeedbackRating: f,
		repoNotif:          g,
		contextTimeOut:     timeout}
}

func (u *useOrder) GetDataBy(ctx context.Context, Claims util.Claims, ID int, GeoUser models.GeoBarber) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		result = models.OrderDataBy{}
	)

	dataHeader, err := u.repoOrderH.GetDataBy(ID, GeoUser)
	if err != nil {
		return nil, err
	}

	dataDetail, err := u.repoOrderD.GetDataBy(ID)
	if err != nil {
		return nil, err
	}

	dataFeedback, err := u.repoFeedbackRating.GetDataBy(ID)
	if err != nil && err != models.ErrNotFound {
		return nil, err
	}

	result.OrderDGet = dataHeader
	result.DataDetail = dataDetail
	result.DataFeedbackRating = dataFeedback

	// dataHeader.DataDetail = dataDetail
	// dataHeader.DataFeedbackRating = dataFeedback

	return result, nil
}
func (u *useOrder) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		queryparam.InitSearch += fmt.Sprintf(" AND user_id = %s", Claims.UserID)
	} else {
		queryparam.InitSearch = fmt.Sprintf("user_id = %s", Claims.UserID)
	}

	UserID, _ := strconv.Atoi(Claims.UserID)
	result.Data, err = u.repoOrderH.GetList(UserID, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoOrderH.Count(UserID, queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useOrder) Create(ctx context.Context, Claims util.Claims, data *models.OrderPost) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mOrder              models.OrderH
		ParamBookingCapster = &models.AddBookingCapster{}
		AddNotif            = &models.AddNotification{}
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mOrder)
	if err != nil {
		return err
	}
	fn := &repofunction.FN{
		Claims: Claims,
	}

	dataBarber, err := fn.GetBarberData(mOrder.BarberID)
	if err != nil {
		if err == models.ErrNotFound {
			return errors.New("Barber sudah tidak bekerja sama.")
		}
		return err
	}
	if !dataBarber.IsActive {
		return errors.New("Tidak bisa order,Barber sedang tidak aktif")
	}

	if !fn.InTimeActiveBarber(dataBarber, data.OrderDate) {
		return errors.New("Mohon maaf , waktu di luar jam oprasional")
	}

	dataCapster, err := fn.GetCapsterData(mOrder.CapsterID)
	if err != nil {
		if err == models.ErrNotFound {
			return errors.New("Data Capster tidak ditemukan, Hubungi pemilik barber")
		}
		return err
	}

	if !dataCapster.IsActive {
		return errors.New("Tidak bisa order,Capster sedang tidak aktif")
	}

	dataUser, err := fn.GetUserData()
	if err != nil {
		return err
	}
	mOrder.OrderDate = data.OrderDate
	mOrder.Status = "N"
	mOrder.FromApps = true
	mOrder.UserID = dataUser.UserID
	mOrder.CustomerName = dataUser.Name
	mOrder.Telp = dataUser.Telp
	mOrder.OrderNo, err = fn.GenTransactionNo(dataBarber.BarberCd)
	if err != nil {
		return err
	}

	mOrder.UserInput = Claims.UserID
	mOrder.UserEdit = Claims.UserID

	//validasi order/booking
	ParamBookingCapster.BarberID = data.BarberID
	ParamBookingCapster.BookingDate = data.OrderDate
	ParamBookingCapster.CapsterID = data.CapsterID

	CntJadwal, _ := u.repoBookingCapster.Count(*ParamBookingCapster, dataUser.UserID)
	if CntJadwal > 0 {
		return errors.New("Mohon Cancel Order/Booking sebelumnya")
	}

	// create order
	err = u.repoOrderH.Create(&mOrder)
	if err != nil {
		return err
	}

	for _, dataDetail := range data.Pakets {
		var mDetail models.OrderD
		err = mapstructure.Decode(dataDetail, &mDetail)
		if err != nil {
			return err
		}
		mDetail.BarberID = mOrder.BarberID
		mDetail.OrderID = mOrder.OrderID
		mDetail.UserEdit = Claims.UserID
		mDetail.UserInput = Claims.UserID
		err = u.repoOrderD.Create(&mDetail)
		if err != nil {
			return err
		}
	}

	//send notif to capster
	//token capster
	capsterFCM := fmt.Sprintf("%v", redisdb.GetSession(strconv.Itoa(mOrder.CapsterID)+"_fcm"))
	barberFCM := fmt.Sprintf("%v", redisdb.GetSession(strconv.Itoa(mOrder.BarberID)+"_fcm"))

	AddNotif.Title = "Ada Orderan"
	// AddNotif.Descs = dataUser.Name + " akan cukur pada :" + mOrder.OrderDate.Format("02/01/2006 15:04")
	AddNotif.Descs = "Ada jadwal cukur pada :" + mOrder.OrderDate.Format("02/01/2006 15:04")
	AddNotif.NotificationStatus = "N"
	AddNotif.NotificationType = "I" // I = Info ; O = Order
	AddNotif.LinkId = mOrder.OrderID
	AddNotif.NotificationDate = util.GetTimeNow()

	if capsterFCM != "" {
		AddNotif.UserId = mOrder.CapsterID
		go u.usenotification.Create(ctx, Claims, capsterFCM, AddNotif)
	}

	if barberFCM != "" {
		AddNotif.UserId = mOrder.BarberID
		go u.usenotification.Create(ctx, Claims, barberFCM, AddNotif)
	}

	// if err != nil {
	// 	return err
	// }

	return nil

}
func (u *useOrder) Update(ctx context.Context, Claims util.Claims, ID int, data models.OrderStatus) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var dataUpdate = map[string]interface{}{
		"status": data.Status,
	}

	err = u.repoOrderH.Update(ID, dataUpdate)
	if err != nil {
		return err
	}

	//delete then insert detail
	if data.Status == "C" {
		var notifUpdate = map[string]interface{}{
			"notification_status": "C",
		}
		err = u.repoNotif.Update(ID, notifUpdate)
		if err != nil {
			return err
		}

	} else {
		err = u.repoOrderD.Delete(ID)
		if err != nil {
			return err
		}

		// for _, dataDetail := range data.Pakets {
		// 	var mDetail models.OrderD
		// 	err = mapstructure.Decode(dataDetail, &mDetail)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	mDetail.BarberID = mOrder.BarberID
		// 	mDetail.OrderID = mOrder.OrderID
		// 	mDetail.UserEdit = Claims.UserID
		// 	mDetail.UserInput = Claims.UserID
		// 	err = u.repoOrderD.Create(&mDetail)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}

	return nil
}
func (u *useOrder) Delete(ctx context.Context, Claims util.Claims, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoOrderH.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
