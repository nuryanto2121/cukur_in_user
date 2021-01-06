package usercapster

import (
	"context"
	"fmt"
	"math"
	ibarbercapster "nuryanto2121/cukur_in_user/interface/barber_capster"
	icapster "nuryanto2121/cukur_in_user/interface/capster"
	ifileupload "nuryanto2121/cukur_in_user/interface/fileupload"
	iuser "nuryanto2121/cukur_in_user/interface/user"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strings"
	"time"
)

type useCapster struct {
	repoCapster       icapster.Repository
	repoUser          iuser.Repository
	repoBarberCapster ibarbercapster.Repository
	repoFile          ifileupload.Repository
	contextTimeOut    time.Duration
}

func NewUserMCapster(a icapster.Repository, b iuser.Repository, c ibarbercapster.Repository, d ifileupload.Repository, timeout time.Duration) icapster.Usecase {
	return &useCapster{
		repoCapster:       a,
		repoUser:          b,
		repoBarberCapster: c,
		repoFile:          d,
		contextTimeOut:    timeout,
	}
}

func (u *useCapster) GetDataBy(ctx context.Context, Claims util.Claims, ID int, GeoBarber models.GeoBarber) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	dataCapster, err := u.repoCapster.GetDataBy(ID, GeoBarber)
	if err != nil {
		return result, err
	}

	dataFile, err := u.repoFile.GetBySaFileUpload(ctx, dataCapster.FileID)
	if err != nil {
		if err != models.ErrNotFound {
			return result, err
		}
	}

	dataCollection, err := u.repoCapster.GetListFileCapter(ID)
	if err != nil {
		if err != models.ErrNotFound {
			return result, err
		}

	}
	response := map[string]interface{}{
		"capster_id":     dataCapster.CapsterID,
		"capster_name":   dataCapster.CapsterName,
		"capster_rating": dataCapster.CapsterRating,
		"join_date":      dataCapster.JoinDate,
		"is_active":      dataCapster.IsActive,
		"is_busy":        dataCapster.IsBusy,
		"file_id":        dataCapster.FileID,
		"file_name":      dataFile.FileName,
		"file_path":      dataFile.FilePath,
		"top_collection": dataCollection,
		"barber_id":      dataCapster.BarberID,
		"barber_name":    dataCapster.BarberName,
		"barber_rating":  dataCapster.BarberRating,
		"distance":       dataCapster.Distance,
		"length_of_work": dataCapster.LengthOfWork,
	}

	return response, nil
}
func (u *useCapster) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		queryparam.InitSearch += fmt.Sprintf(` 
		AND distance <= 10
		and is_barber_active = true
		and is_barber_open = true
		and is_active = true`)
	} else {
		queryparam.InitSearch = fmt.Sprintf(`distance <= 10 and is_barber_active = true
		and is_barber_open = true
		and is_active = true`)
	}
	result.Data, err = u.repoBarberCapster.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoBarberCapster.Count(queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
