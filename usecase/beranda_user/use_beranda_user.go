package useberandauser

import (
	"context"
	"fmt"
	"math"
	iadvertis "nuryanto2121/cukur_in_user/interface/advertise"
	ibarber "nuryanto2121/cukur_in_user/interface/barber"
	ibarbercapster "nuryanto2121/cukur_in_user/interface/barber_capster"
	iberandauser "nuryanto2121/cukur_in_user/interface/beranda_user"
	ifileupload "nuryanto2121/cukur_in_user/interface/fileupload"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strconv"
	"strings"
	"time"
)

type useBerandaUser struct {
	useBarber         ibarber.Usecase
	repoFile          ifileupload.Repository
	repoBarberCapster ibarbercapster.Repository
	repoBarber        ibarber.Repository
	useAdvertise      iadvertis.Usecase
	contextTimeOut    time.Duration
}

func NewUseBerandaUser(a ibarber.Usecase, b ifileupload.Repository, c ibarbercapster.Repository, d ibarber.Repository, e iadvertis.Usecase, timeout time.Duration) iberandauser.Usecase {
	return &useBerandaUser{
		useBarber:         a,
		repoFile:          b,
		repoBarberCapster: c,
		repoBarber:        d,
		useAdvertise:      e,
		contextTimeOut:    timeout,
	}
}

func (u *useBerandaUser) GetClosestBarber(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	queryparam.PerPage = 5
	queryparam.SortField = "distance ASC"
	queryparam.InitSearch = "is_active = true AND distance <= 10" //is_barber_open = true
	// queryparam.InitSearch = "is_active = 't' AND distance <= 10 "
	// result, err = u.useBarber.GetList(ctx, Claims, queryparam)
	// if err != nil {
	// 	return result, err
	// }
	ID, _ := strconv.Atoi(Claims.UserID)
	result.Data, err = u.repoBarber.GetList(ID, queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoBarber.Count(ID, queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useBerandaUser) GetRecomentCapster(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	queryparam.PerPage = 7
	queryparam.SortField = "capster_rating desc,distance"
	if queryparam.Search != "" {
		// queryparam.Search = fmt.Sprintf("lower(barber_name) iLIKE '%%%s%%' ", queryparam.Search)
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" { //		and is_barber_open = true
		queryparam.InitSearch += fmt.Sprintf(` 
		AND distance <= 10
		and is_barber_active = true
		and is_active = true`)
	} else {
		queryparam.InitSearch = fmt.Sprintf(`
		distance <= 10 
		and is_barber_active = true		
		and is_active = true`) //and is_barber_open = true
	}
	result.Data, err = u.repoBarberCapster.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoBarberCapster.Count(queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useBerandaUser) GetRecomentBarber(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	queryparam.PerPage = 5
	queryparam.SortField = "barber_rating DESC,distance"
	// queryparam.InitSearch = "is_active = true and now()::time >= b.operation_start::time  and now()::time <= b.operation_end::time" // (now()::time between b.operation_start::time and b.operation_end::time )"
	queryparam.InitSearch = "is_active = true " //and is_barber_open = true
	result, err = u.useBarber.GetList(ctx, Claims, queryparam)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *useBerandaUser) GetAdvertiseBarber(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		queryparamx models.ParamList
	)
	queryparamx.PerPage = 5
	queryparamx.SortField = "advertise_id"
	result, err = u.useAdvertise.GetList(ctx, Claims, queryparamx)
	if err != nil {
		return result, err
	}

	return result, nil
}
