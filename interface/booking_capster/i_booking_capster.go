package ibookingcapster

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Repository interface {
	GetDataBy(queryparam models.AddBookingCapster) (result *[]models.OrderH, err error)
	Count(param models.AddBookingCapster, UserID int) (result int, err error)
	Create(data *models.BookingCapster) (err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, param models.AddBookingCapster) (result interface{}, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.AddBookingCapster) error
}
