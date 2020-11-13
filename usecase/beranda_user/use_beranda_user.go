package useberandauser

import (
	"context"
	ibarber "nuryanto2121/cukur_in_user/interface/barber"
	iberandauser "nuryanto2121/cukur_in_user/interface/beranda_user"
	ifileupload "nuryanto2121/cukur_in_user/interface/fileupload"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"time"
)

type useBerandaUser struct {
	useBarber      ibarber.Usecase
	repoFile       ifileupload.Repository
	contextTimeOut time.Duration
}

func NewUseBerandaUser(a ibarber.Usecase, b ifileupload.Repository, timeout time.Duration) iberandauser.Usecase {
	return &useBerandaUser{
		useBarber:      a,
		repoFile:       b,
		contextTimeOut: timeout,
	}
}

func (u *useBerandaUser) GetClosestBarber(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	queryparam.PerPage = 5
	queryparam.SortField = "distance"
	result, err = u.useBarber.GetList(ctx, Claims, queryparam)
	if err != nil {
		return result, err
	}

	return result, nil
}
func (u *useBerandaUser) GetRecomentCapster(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error)
func (u *useBerandaUser) GetRecomentBarber(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	queryparam.PerPage = 5
	queryparam.SortField = "barber_rating,distance"
	result, err = u.useBarber.GetList(ctx, Claims, queryparam)
	if err != nil {
		return result, err
	}

	return result, nil
}
