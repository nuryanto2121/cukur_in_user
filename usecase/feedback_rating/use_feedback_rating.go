package usefeedbackrating

import (
	"context"
	"fmt"
	"math"
	ifeedbackrating "nuryanto2121/cukur_in_user/interface/feedback_rating"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type useFeedbackRating struct {
	repoFeedbackRating ifeedbackrating.Repository
	contextTimeOut     time.Duration
}

func NewFeedbackRating(a ifeedbackrating.Repository, timeout time.Duration) ifeedbackrating.Usecase {
	return &useFeedbackRating{repoFeedbackRating: a, contextTimeOut: timeout}
}
func (u *useFeedbackRating) GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	result, err = u.repoFeedbackRating.GetDataBy(ID)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (u *useFeedbackRating) Update(ctx context.Context, Claims util.Claims, ID int, data models.AddFeedbackRating) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	datas := structs.Map(data)
	datas["user_edit"] = Claims.UserID
	err = u.repoFeedbackRating.Update(ID, datas)
	if err != nil {
		return err
	}
	return nil
}
func (u *useFeedbackRating) Create(ctx context.Context, Claims util.Claims, data *models.AddFeedbackRating) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mPaket models.FeedbackRating
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mPaket.AddFeedbackRating)
	if err != nil {
		return err
	}
	mPaket.UserID, _ = strconv.Atoi(Claims.UserID)
	mPaket.UserEdit = Claims.UserID
	mPaket.UserInput = Claims.UserID

	err = u.repoFeedbackRating.Create(&mPaket)
	if err != nil {
		return err
	}
	return nil
}
func (u *useFeedbackRating) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	// var tUser = models.Paket{}
	/*membuat Where like dari struct*/
	if queryparam.Search != "" {
		// queryparam.Search = fmt.Sprintf("paket_name iLIKE '%%%s%%' OR descs iLIKE '%%%s%%'", queryparam.Search, queryparam.Search)
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		queryparam.InitSearch += " AND user_id = " + Claims.UserID
	} else {
		queryparam.InitSearch = " user_id = " + Claims.UserID
	}
	result.Data, err = u.repoFeedbackRating.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoFeedbackRating.Count(queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
