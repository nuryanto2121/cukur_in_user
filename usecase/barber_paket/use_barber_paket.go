package usebarberpaket

import (
	"context"
	ibarberpaket "nuryanto2121/cukur_in_user/interface/barber_paket"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"time"
)

type useBarberPaket struct {
	repoBarberPaket ibarberpaket.Repository
	contextTimeOut  time.Duration
}

func NewBarberPaket(a ibarberpaket.Repository, timeout time.Duration) ibarberpaket.Usecase {
	return &useBarberPaket{repoBarberPaket: a, contextTimeOut: timeout}
}

func (u *useBarberPaket) GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		queryparam models.ParamListGeo
	)

	dataBPaket, err := u.repoBarberPaket.GetList(queryparam)
	if err != nil {
		return result, err
	}
	return dataBPaket, nil
}
