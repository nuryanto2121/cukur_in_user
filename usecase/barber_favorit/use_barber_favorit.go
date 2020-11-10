package usebarberfavorit

import (
	"context"
	"fmt"
	"math"
	ibarberfavorit "nuryanto2121/cukur_in_user/interface/barber_favorit"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

type useBarberFavorit struct {
	repoBarberFavorit ibarberfavorit.Repository
	contextTimeOut    time.Duration
}

func NewBarberFavorit(a ibarberfavorit.Repository, timeout time.Duration) ibarberfavorit.Usecase {
	return &useBarberFavorit{repoBarberFavorit: a, contextTimeOut: timeout}
}
func (u *useBarberFavorit) Create(ctx context.Context, Claims util.Claims, data *models.AddBarberFavorit) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mBarberFavorit models.BarberFavorit
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mBarberFavorit)
	if err != nil {
		return err
	}
	mBarberFavorit.UserID, _ = strconv.Atoi(Claims.UserID)
	mBarberFavorit.UserEdit = Claims.UserID
	mBarberFavorit.UserInput = Claims.UserID

	_, err = u.repoBarberFavorit.GetDataBy(data.BarberID, mBarberFavorit.UserID)
	if err != nil && err == models.ErrNotFound {
		err = u.repoBarberFavorit.Create(&mBarberFavorit)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}

		err = u.repoBarberFavorit.Delete(data.BarberID, mBarberFavorit.UserID)

	}

	return nil
}
func (u *useBarberFavorit) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	// var tUser = models.BarberFavorit{}
	if queryparam.Search != "" {
		// queryparam.Search = fmt.Sprintf("paket_name iLIKE '%%%s%%' OR descs iLIKE '%%%s%%'", queryparam.Search, queryparam.Search)
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		queryparam.InitSearch += " AND a.user_id = " + Claims.UserID
	} else {
		queryparam.InitSearch = " a.user_id = " + Claims.UserID
	}
	result.Data, err = u.repoBarberFavorit.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoBarberFavorit.Count(queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

// func (u *useBarberFavorit) Delete(ctx context.Context, ID int) (err error) {
// 	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
// 	defer cancel()

// 	err = u.repoBarberFavorit.Delete(ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
