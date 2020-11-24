package usebookingcapster

import (
	"context"
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

	response := map[string]interface{}{
		"list_booking_capster":  JadwalCapster,
		"list_hour_operational": "",
	}

	return response, nil
}
func (u *useBookingCapster) Create(ctx context.Context, Claims util.Claims, data *models.AddBookingCapster) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mBookingCapster models.BookingCapster
	)

	// mapping to struct model saRole
	err := mapstructure.Decode(data, &mBookingCapster)
	if err != nil {
		return err
	}
	mBookingCapster.UserEdit = Claims.UserID
	mBookingCapster.UserInput = Claims.UserID

	err = u.repoBookingCapster.Create(&mBookingCapster)
	if err != nil {
		return err
	}
	return nil
}
