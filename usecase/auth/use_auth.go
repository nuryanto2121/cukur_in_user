package useauth

import (
	"context"
	"errors"
	iauth "nuryanto2121/cukur_in_user/interface/auth"
	ifileupload "nuryanto2121/cukur_in_user/interface/fileupload"
	iusers "nuryanto2121/cukur_in_user/interface/user"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/setting"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"nuryanto2121/cukur_in_user/redisdb"
	useemailauth "nuryanto2121/cukur_in_user/usecase/email_auth"
	"strconv"
	"time"
)

type useAuht struct {
	repoAuth       iusers.Repository
	repoFile       ifileupload.Repository
	contextTimeOut time.Duration
}

func NewUserAuth(a iusers.Repository, b ifileupload.Repository, timeout time.Duration) iauth.Usecase {
	return &useAuht{repoAuth: a, repoFile: b, contextTimeOut: timeout}
}
func (u *useAuht) Login(ctx context.Context, dataLogin *models.LoginForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		expireToken = setting.FileConfigSetting.JWTExpired
	)

	DataUser, err := u.repoAuth.GetByAccount(dataLogin.Account) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
	if err != nil {
		// return util.GoutputErrCode(http.StatusUnauthorized, "Your User/Email not valid.") //appE.ResponseError(util.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
		return nil, errors.New("Your Account not valid.")
	}

	if !util.ComparePassword(DataUser.Password, util.GetPassword(dataLogin.Password)) {
		return nil, errors.New("Your Password not valid.")
	}
	DataFile, err := u.repoFile.GetBySaFileUpload(ctx, DataUser.FileID)

	token, err := util.GenerateToken(DataUser.UserID, dataLogin.Account, DataUser.UserType)
	if err != nil {
		return nil, err
	}

	redisdb.AddSession(token, DataUser.UserID, time.Duration(expireToken)*time.Hour)

	restUser := map[string]interface{}{
		"user_id":   DataUser.UserID,
		"email":     DataUser.Email,
		"telp":      DataUser.Telp,
		"user_name": DataUser.Name,
		"user_type": DataUser.UserType,
		"file_id":   DataUser.FileID,
		"file_name": DataFile.FileName,
		"file_path": DataFile.FilePath,
	}
	response := map[string]interface{}{
		"token":     token,
		"data_user": restUser,
	}

	return response, nil
}

func (u *useAuht) ForgotPassword(ctx context.Context, dataForgot *models.ForgotForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	DataUser, err := u.repoAuth.GetByAccount(dataForgot.Account) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
	if err != nil {
		// return util.GoutputErrCode(http.StatusUnauthorized, "Your User/Email not valid.") //appE.ResponseError(util.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
		return errors.New("Your Account not valid.")
	}
	if DataUser.Name == "" {
		return errors.New("Your Account not valid.")
	}

	return nil
}

func (u *useAuht) ResetPassword(ctx context.Context, dataReset *models.ResetPasswd) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if dataReset.Passwd != dataReset.ConfirmPasswd {
		return errors.New("Password and Confirm Password not same.")
	}

	DataUser, err := u.repoAuth.GetByAccount(dataReset.Account)
	if err != nil {
		return err
	}

	DataUser.Password, _ = util.Hash(dataReset.Passwd)
	// email, err := util.ParseEmailToken(dataReset.TokenEmail)
	// if err != nil {
	// 	email = dataReset.TokenEmail
	// }

	// dataUser, err := u.repoUser.GetByEmailSaUser(email)

	// dataUser.Password, _ = util.Hash(dataReset.Passwd)

	err = u.repoAuth.Update(DataUser.UserID, &DataUser)
	if err != nil {
		return err
	}

	return nil
}

func (u *useAuht) Register(ctx context.Context, dataRegister models.RegisterForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var User models.SsUser

	CekData, err := u.repoAuth.GetByAccount(dataRegister.Account) //repoAuth.GetByAccount(dataRegister.EmailAddr, true)
	if CekData.Email == dataRegister.Account {
		return output, errors.New("email sudah terdaftar.")
	}

	if dataRegister.Passwd != dataRegister.ConfirmPasswd {
		return output, errors.New("Password dan confirm password tidak boleh sama.")
	}
	User.Name = dataRegister.Name
	User.Password, _ = util.Hash(dataRegister.Passwd)
	User.UserEdit = dataRegister.Name
	User.UserInput = dataRegister.Name
	User.JoinDate = time.Now()
	User.UserType = "user"
	//check email or telp
	// if util.CheckEmail(dataRegister.Account) {
	User.Email = dataRegister.Account
	// } else {
	// User.Telp = dataRegister.Account
	// }
	err = u.repoAuth.Create(&User)
	if err != nil {
		return output, err
	}

	mUser := map[string]interface{}{
		"user_input": strconv.Itoa(User.UserID),
		"user_edit":  strconv.Itoa(User.UserID),
	}
	err = u.repoAuth.Update(User.UserID, mUser)
	if err != nil {
		return output, err
	}

	GenCode := util.GenerateNumber(4)

	// send generate code
	mailService := &useemailauth.Register{
		Email:      User.Email,
		Name:       User.Name,
		PasswordCd: GenCode,
	}

	go mailService.SendRegister()
	// err = mailService.SendRegister()
	// if err != nil {
	// 	return output, err
	// }

	//store to redis
	err = redisdb.AddSession(dataRegister.Account, GenCode, 24*time.Hour)
	if err != nil {
		return output, err
	}
	out := map[string]interface{}{
		"otp":     GenCode,
		"account": User.Email,
	}
	return out, nil
}

func (u *useAuht) Verify(ctx context.Context, dataVeriry models.VerifyForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()
	data := redisdb.GetSession(dataVeriry.Account)
	if data == "" {
		return errors.New("Please Resend Code")
	}

	if data != dataVeriry.VerifyCode {
		return errors.New("Invalid Code.")
	}

	return nil
}
