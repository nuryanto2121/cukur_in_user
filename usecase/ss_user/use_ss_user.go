package usesysuser

import (
	"context"
	"errors"
	"math"
	ifileupload "nuryanto2121/cukur_in_user/interface/fileupload"
	iusers "nuryanto2121/cukur_in_user/interface/user"
	"nuryanto2121/cukur_in_user/models"
	querywhere "nuryanto2121/cukur_in_user/pkg/query"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"reflect"
	"strconv"
	"time"

	"github.com/fatih/structs"
)

type useSysUser struct {
	repoUser       iusers.Repository
	repoFile       ifileupload.Repository
	contextTimeOut time.Duration
}

func NewUserSysUser(a iusers.Repository, b ifileupload.Repository, timeout time.Duration) iusers.Usecase {
	return &useSysUser{repoUser: a, repoFile: b, contextTimeOut: timeout}
}

func (u *useSysUser) GetByEmailSaUser(ctx context.Context, email string) (result models.SsUser, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	a := models.SsUser{}
	result, err = u.repoUser.GetByAccount(email)
	if err != nil {
		return a, err
	}
	return result, nil
}
func (u *useSysUser) ChangePassword(ctx context.Context, Claims util.Claims, DataChangePwd models.ChangePassword) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	ID, _ := strconv.Atoi(Claims.UserID)
	dataUser, err := u.repoUser.GetDataBy(ID)
	if err != nil {
		return err
	}

	// DataChangePwd.OldPassword, _ = util.Hash(DataChangePwd.OldPassword)
	if !util.ComparePassword(dataUser.Password, util.GetPassword(DataChangePwd.OldPassword)) {
		return errors.New("Password lama anda salah.")
	}

	if DataChangePwd.NewPassword != DataChangePwd.ConfirmPassword {
		return errors.New("Password dan confirm password tidak sama.")
	}

	if util.ComparePassword(dataUser.Password, util.GetPassword(DataChangePwd.NewPassword)) {
		return errors.New("Password baru tidak boleh sama dengan yang lama.")
	}

	DataChangePwd.NewPassword, _ = util.Hash(DataChangePwd.NewPassword)
	// var data = map[string]interface{}{
	// 	"password": DataChangePwd.NewPassword,
	// }

	// err = u.repoUser.Update(ID, data)
	err = u.repoUser.UpdatePasswordByEmail(dataUser.Email, DataChangePwd.NewPassword)
	if err != nil {
		return err
	}
	return nil
}
func (u *useSysUser) GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	DataUser, err := u.repoUser.GetDataBy(ID)
	if err != nil {
		return result, err
	}
	DataFile, err := u.repoFile.GetBySaFileUpload(ctx, DataUser.FileID)
	if err != nil {
		if err != models.ErrNotFound {
			return result, err
		}
	}
	response := map[string]interface{}{
		"user_id":       DataUser.UserID,
		"user_name":     DataUser.Name,
		"birth_of_date": DataUser.BirthOfDate,
		"email":         DataUser.Email,
		"telp":          DataUser.Telp,
		"file_id":       DataUser.FileID,
		"file_name":     DataFile.FileName,
		"file_path":     DataFile.FilePath,
	}
	return response, nil
}
func (u *useSysUser) GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var tUser = models.SsUser{}
	/*membuat Where like dari struct*/
	if queryparam.Search != "" {
		value := reflect.ValueOf(tUser)
		types := reflect.TypeOf(&tUser)
		queryparam.Search = querywhere.GetWhereLikeStruct(value, types, queryparam.Search, "")
	}
	result.Data, err = u.repoUser.GetList(queryparam)
	if err != nil {
		return result, err
	}

	result.Total, err = u.repoUser.Count(queryparam)
	if err != nil {
		return result, err
	}

	// d := float64(result.Total) / float64(queryparam.PerPage)
	result.LastPage = int(math.Ceil(float64(result.Total) / float64(queryparam.PerPage)))
	result.Page = queryparam.Page

	return result, nil
}
func (u *useSysUser) Create(ctx context.Context, Claims util.Claims, data *models.SsUser) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUser.Create(data)
	if err != nil {
		return err
	}
	return nil

}
func (u *useSysUser) Update(ctx context.Context, Claims util.Claims, ID int, data models.UpdateUser) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	dataUser, err := u.repoUser.GetByAccount(data.Email)
	if dataUser.UserID != ID {
		return errors.New("Email sudah terdaftar.")
	}

	datas := structs.Map(data)
	datas["user_edit"] = Claims.UserID
	err = u.repoUser.Update(ID, datas)
	if err != nil {
		return err
	}
	return nil
}
func (u *useSysUser) Delete(ctx context.Context, ID int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	err = u.repoUser.Delete(ID)
	if err != nil {
		return err
	}
	return nil
}
