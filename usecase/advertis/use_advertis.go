package useadvertis

import (
	"context"
	"fmt"
	"math"
	iadvertis "nuryanto2121/cukur_in_user/interface/advertis"
	ifileupload "nuryanto2121/cukur_in_user/interface/fileupload"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type useAdvertis struct {
	repoAdvertis   iadvertis.Repository
	repoFile       ifileupload.Repository
	contextTimeOut time.Duration
}

func NewUseAdvertis(a iadvertis.Repository, b ifileupload.Repository, timeout time.Duration) iadvertis.Usecase {
	return &useAdvertis{
		repoAdvertis:   a,
		repoFile:       b,
		contextTimeOut: timeout,
	}
}

func (u *useAdvertis) GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	dataAdvertis, err := u.repoAdvertis.GetDataBy(ID)
	if err != nil {
		return result, err
	}
	dataFile, err := u.repoFile.GetBySaFileUpload(ctx, dataAdvertis.FileID)
	if err != nil {
		return result, err
	}

	result = map[string]interface{}{
		"advertis_id":     dataAdvertis.AdvertisID,
		"title":           dataAdvertis.Title,
		"descs":           dataAdvertis.Descs,
		"advertis_status": dataAdvertis.AdvertisStatus,
		"slide_duration":  dataAdvertis.SlideDuration,
		"start_date":      dataAdvertis.StartDate,
		"end_date":        dataAdvertis.EndDate,
		"file_id":         dataFile.FileID,
		"file_name":       dataFile.FileName,
		"file_path":       dataFile.FilePath,
		"file_type":       dataFile.FileType,
	}
	return result, nil
}

func (u *useAdvertis) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
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
	result.Data, err = u.repoAdvertis.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoAdvertis.Count(queryparam)
	if err != nil {
		return result, err
	}

	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}

func (u *useAdvertis) Create(ctx context.Context, Claims util.Claims, data *models.AddAdvertis) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mAdvertis = models.Advertis{}
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mAdvertis.AddAdvertis)
	if err != nil {
		return err
	}

	mAdvertis.UserEdit = Claims.UserID
	mAdvertis.UserInput = Claims.UserID

	err = u.repoAdvertis.Create(&mAdvertis)
	if err != nil {
		return err
	}
	return nil

}

func (u *useAdvertis) Update(ctx context.Context, Claims util.Claims, ID int, data *models.AddAdvertis) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	myMap := structs.Map(data)
	myMap["user_edit"] = Claims.UserID
	fmt.Println(myMap)
	err = u.repoAdvertis.Update(ID, myMap)
	if err != nil {
		return err
	}
	return nil
}

func (u *useAdvertis) Delete(ctx context.Context, Claims util.Claims, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoAdvertis.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
