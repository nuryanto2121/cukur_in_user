package usebookingcapster

import (
	"context"
	"errors"
	ibookingcapster "nuryanto2121/cukur_in_user/interface/booking_capster"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"time"

	"github.com/mitchellh/mapstructure"
)

type useBookingCapster struct {
	repoBookingCapster ibookingcapster.Repository
	contextTimeOut     time.Duration
}

func NewBookingCapster(a ibookingcapster.Repository, timeout time.Duration) ibookingcapster.Usecase {
	return &useBookingCapster{repoBookingCapster: a, contextTimeOut: timeout}
}

func (u *useBookingCapster) GetDataBy(ctx context.Context, Claims util.Claims, queryparam models.AddBookingCapster) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	JadwalCapster, err := u.repoBookingCapster.GetDataBy(queryparam)
	if err != nil {
		return result, err
	}

	// response := map[string]interface{}{
	// 	"list_booking_capster":  JadwalCapster,
	// 	"list_hour_operational": "",
	// }

	return JadwalCapster, nil
}
func (u *useBookingCapster) Create(ctx context.Context, Claims util.Claims, data *models.AddBookingCapster) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mBookingCapster models.BookingCapster
	)

	// mapping to struct model saRole
	err := mapstructure.Decode(data, &mBookingCapster.AddBookingCapster)
	if err != nil {
		return err
	}
	mBookingCapster.UserEdit = Claims.UserID
	mBookingCapster.UserInput = Claims.UserID

	Cnt, err := u.repoBookingCapster.Count(mBookingCapster.AddBookingCapster, 0)
	if Cnt > 0 {
		return errors.New("Waktu sudah ada yang booking.")
	}
	// if err != nil {
	// 	return  err
	// }

	err = u.repoBookingCapster.Create(&mBookingCapster)
	if err != nil {
		return err
	}
	return nil
}
