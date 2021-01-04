package usebarber

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	repofunction "nuryanto2121/cukur_in_user/repository/function"

	ibarber "nuryanto2121/cukur_in_user/interface/barber"
	ibarbercapster "nuryanto2121/cukur_in_user/interface/barber_capster"
	ibarberpaket "nuryanto2121/cukur_in_user/interface/barber_paket"
	ifileupload "nuryanto2121/cukur_in_user/interface/fileupload"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type useBarber struct {
	repoBarber        ibarber.Repository
	repoBarberPaket   ibarberpaket.Repository
	repoBarberCapster ibarbercapster.Repository
	repoFile          ifileupload.Repository
	contextTimeOut    time.Duration
}

func NewUserMBarber(a ibarber.Repository, b ibarberpaket.Repository, c ibarbercapster.Repository, d ifileupload.Repository, timeout time.Duration) ibarber.Usecase {
	return &useBarber{
		repoBarber:        a,
		repoBarberPaket:   b,
		repoBarberCapster: c,
		repoFile:          d,
		contextTimeOut:    timeout}
}

func (u *useBarber) GetDataBy(ctx context.Context, Claims util.Claims, ID int, GeoBarber models.GeoBarber) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		queryparam models.ParamListGeo
	)
	queryparam.Latitude = GeoBarber.Latitude
	queryparam.Longitude = GeoBarber.Longitude

	UserID, _ := strconv.Atoi(Claims.UserID)
	result, err := u.repoBarber.GetDataByList(ID, UserID, GeoBarber)
	if err != nil {
		return result, err
	}
	queryparam.InitSearch = fmt.Sprintf("barber_paket.barber_id = %d", result.BarberID)
	queryparam.Page = 1
	queryparam.PerPage = 50
	dataFile, err := u.repoFile.GetBySaFileUpload(ctx, result.FileID)
	if err != nil {
		return result, err
	}

	dataBPaket, err := u.repoBarberPaket.GetList(queryparam)
	if err != nil {
		return result, err
	}

	queryparam.InitSearch = fmt.Sprintf("barber_id = %d", result.BarberID)
	dataBCapster, err := u.repoBarberCapster.GetList(queryparam)
	if err != nil {
		return result, err
	}
	response := map[string]interface{}{
		"barber_id":        result.BarberID,
		"barber_cd":        result.BarberCd,
		"barber_name":      result.BarberName,
		"address":          result.Address,
		"latitude":         result.Latitude,
		"longitude":        result.Longitude,
		"distance":         result.Distance,
		"barber_rating":    result.BarberRating,
		"is_favorit":       result.IsFavorit,
		"is_barber_open":   result.IsBarberOpen,
		"total_user_order": result.TotalUserOrder,
		"operation_start":  result.OperationStart,
		"operation_end":    result.OperationEnd,
		"is_active":        result.IsActive,
		"file_id":          dataFile.FileID,
		"file_name":        dataFile.FileName,
		"file_path":        dataFile.FilePath,
		"barber_paket":     dataBPaket,
		"barber_capster":   dataBCapster,
	}

	return response, nil
}
func (u *useBarber) GetDataFirst(ctx context.Context, Claims util.Claims, ID int) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	OwnerID, _ := strconv.Atoi(Claims.UserID)
	result, err := u.repoBarber.GetDataFirst(OwnerID, ID)
	if err != nil {
		return result, err
	}

	dataFile, err := u.repoFile.GetBySaFileUpload(ctx, result.FileID)
	if err != nil {
		return result, err
	}

	response := map[string]interface{}{
		"barber_id":       result.BarberID,
		"barber_cd":       result.BarberCd,
		"barber_name":     result.BarberName,
		"address":         result.Address,
		"latitude":        result.Latitude,
		"longitude":       result.Longitude,
		"operation_start": result.OperationStart,
		"operation_end":   result.OperationEnd,
		"is_active":       result.IsActive,
		"file_id":         dataFile.FileID,
		"file_name":       dataFile.FileName,
		"file_path":       dataFile.FilePath,
	}

	return response, nil
}
func (u *useBarber) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	// var tUser = models.Barber{}
	/*membuat Where like dari struct*/
	if queryparam.Search != "" {
		// queryparam.Search = fmt.Sprintf("lower(barber_name) iLIKE '%%%s%%' ", queryparam.Search)
		queryparam.Search = strings.ToLower(fmt.Sprintf("%%%s%%", queryparam.Search))
	}

	if queryparam.InitSearch != "" {
		queryparam.InitSearch += fmt.Sprintf(" AND distance <= 10")
	} else {
		queryparam.InitSearch = fmt.Sprintf("distance <= 10")
	}
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
func (u *useBarber) Create(ctx context.Context, Claims util.Claims, data *models.BarbersPost) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	var (
		mBarber models.Barber
		// fn      = &repofunction.FN{Claims}
	)
	fn := &repofunction.FN{
		Claims: Claims,
	}
	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mBarber)
	if err != nil {
		return err
	}
	mBarber.BarberCd, err = fn.GenBarberCode()
	if err != nil {
		return err
	}
	mBarber.OperationEnd = data.OperationEnd
	mBarber.OperationStart = data.OperationStart
	mBarber.OwnerID, _ = strconv.Atoi(Claims.UserID)
	mBarber.UserInput = Claims.UserID
	mBarber.UserEdit = Claims.UserID
	err = u.repoBarber.Create(&mBarber)
	if err != nil {
		return err
	}

	for _, dataCapster := range data.BarberCapster {
		var BCapster = models.BarberCapster{}
		BCapster.BarberID = mBarber.BarberID
		BCapster.CapsterID = dataCapster.CapsterID
		BCapster.UserInput = Claims.UserID
		BCapster.UserEdit = Claims.UserID
		err = u.repoBarberCapster.Create(&BCapster)
		if err != nil {
			return err
		}
	}

	for _, dataCapster := range data.BarberPaket {
		var BPaket = models.BarberPaket{}
		BPaket.BarberID = mBarber.BarberID
		BPaket.PaketID = dataCapster.PaketID
		BPaket.UserInput = Claims.UserID
		BPaket.UserEdit = Claims.UserID
		err = u.repoBarberPaket.Create(&BPaket)
		if err != nil {
			return err
		}
	}

	return nil

}
func (u *useBarber) Update(ctx context.Context, Claims util.Claims, ID int, data models.BarbersPost) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		mBarber models.BarbersUpdate
	)

	// mapping to struct model saRole
	err = mapstructure.Decode(data, &mBarber)
	if err != nil {
		return err
	}

	// check jika status tidak aktif
	if !mBarber.IsActive {
		fn := &repofunction.FN{
			Claims: Claims,
		}
		cnt := fn.GetCountTrxBarber(ID)
		if cnt > 0 {
			return errors.New("Mohon maaf anda tidak dapat non-aktifkan Barber, sedang ada transaksi yang berlangsung")
		}
	}
	mBarber.OperationStart = data.OperationStart
	mBarber.OperationEnd = data.OperationEnd

	dataUpdate := structs.Map(mBarber)
	dataUpdate["user_edit"] = Claims.UserID
	err = u.repoBarber.Update(ID, dataUpdate)
	if err != nil {
		return err
	}

	//delete then insert detail

	err = u.repoBarberCapster.Delete(ID)
	if err != nil {
		return err
	}

	for _, dataCapster := range data.BarberCapster {
		var BCapster = models.BarberCapster{}
		BCapster.BarberID = ID
		BCapster.CapsterID = dataCapster.CapsterID
		BCapster.UserInput = Claims.UserID
		BCapster.UserEdit = Claims.UserID
		err = u.repoBarberCapster.Create(&BCapster)
		if err != nil {
			return err
		}
	}

	err = u.repoBarberPaket.Delete(ID)
	if err != nil {
		return err
	}
	for _, dataCapster := range data.BarberPaket {
		var BPaket = models.BarberPaket{}
		BPaket.BarberID = ID
		BPaket.PaketID = dataCapster.PaketID
		BPaket.UserInput = Claims.UserID
		BPaket.UserEdit = Claims.UserID
		err = u.repoBarberPaket.Create(&BPaket)
		if err != nil {
			return err
		}
	}

	return nil
}
func (u *useBarber) Delete(ctx context.Context, Claims util.Claims, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoBarber.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
