package useadvertise

import (
	"context"
	"fmt"
	"math"
	iadvertis "nuryanto2121/cukur_in_user/interface/advertise"
	ifileupload "nuryanto2121/cukur_in_user/interface/fileupload"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type useAdvertise struct {
	repoAdvertise  iadvertis.Repository
	repoFile       ifileupload.Repository
	contextTimeOut time.Duration
}

func NewUseAdvertise(a iadvertis.Repository, b ifileupload.Repository, timeout time.Duration) iadvertis.Usecase {
	return &useAdvertise{
		repoAdvertise:  a,
		repoFile:       b,
		contextTimeOut: timeout,
	}
}

func (u *useAdvertise) GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	dataAdvertise, err := u.repoAdvertise.GetDataBy(ID)
	if err != nil {
		return result, err
	}
	dataFile, err := u.repoFile.GetBySaFileUpload(ctx, dataAdvertise.FileID)
	if err != nil {
		return result, err
	}

	result = map[string]interface{}{
		"advertise_id":     dataAdvertise.AdvertiseID,
		"title":            dataAdvertise.Title,
		"descs":            dataAdvertise.Descs,
		"advertise_status": dataAdvertise.AdvertiseStatus,
		"slide_duration":   dataAdvertise.SlideDuration,
		"start_date":       dataAdvertise.StartDate,
		"end_date":         dataAdvertise.EndDate,
		"file_id":          dataFile.FileID,
		"file_name":        dataFile.FileName,
		"file_path":        dataFile.FilePath,
		"file_type":        dataFile.FileType,
	}
	return result, nil
}

func (u *useAdvertise) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if queryparam.Search != "" {
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		queryparam.InitSearch += " AND (now()::timestamp between start_date and end_date)"
	} else {
		queryparam.InitSearch = " (now()::timestamp between start_date and end_date)"
	}
	result.Data, err = u.repoAdvertise.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoAdvertise.Count(queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useAdvertise) Create(ctx context.Context, Claims util.Claims, data *models.AddAdvertise) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mAdvertise = models.Advertise{}
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mAdvertise.AddAdvertise)
	if err != nil {
		return err
	}

	mAdvertise.UserEdit = Claims.UserID
	mAdvertise.UserInput = Claims.UserID

	err = u.repoAdvertise.Create(&mAdvertise)
	if err != nil {
		return err
	}
	return nil

}

func (u *useAdvertise) Update(ctx context.Context, Claims util.Claims, ID int, data *models.AddAdvertise) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoAdvertise.Update(ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useAdvertise) Delete(ctx context.Context, Claims util.Claims, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoAdvertise.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
