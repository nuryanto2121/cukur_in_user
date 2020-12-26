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

func (u *useAuht) Logout(ctx context.Context, Claims util.Claims, Token string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	redisdb.TurncateList(Token)
	redisdb.TurncateList(Claims.UserID + "_fcm")

	return nil
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
		return nil, errors.New("Email anda belum terdaftar.")
	}

	if !DataUser.IsActive {
		return nil, errors.New("Account anda belum aktif. Silahkan Register ulang dengan email yang sama.")
	}

	if !util.ComparePassword(DataUser.Password, util.GetPassword(dataLogin.Password)) {
		return nil, errors.New("Password yang anda masukkan salah. Silahkan coba lagi..")
	}
	DataFile, err := u.repoFile.GetBySaFileUpload(ctx, DataUser.FileID)

	token, err := util.GenerateToken(DataUser.UserID, dataLogin.Account, DataUser.UserType)
	if err != nil {
		return nil, err
	}

	redisdb.AddSession(token, DataUser.UserID, time.Duration(expireToken)*time.Hour)

	// expired FCM
	redisdb.AddSession(strconv.Itoa(DataUser.UserID)+"_fcm", dataLogin.FcmToken, time.Duration(expireToken)*time.Hour)

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
		"user_type": "user",
	}

	return response, nil
}

func (u *useAuht) ForgotPassword(ctx context.Context, dataForgot *models.ForgotForm) (result string, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	DataUser, err := u.repoAuth.GetByAccount(dataForgot.Account) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
	if err != nil {
		// return util.GoutputErrCode(http.StatusUnauthorized, "Your User/Email not valid.") //appE.ResponseError(util.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
		return "", errors.New("Your Account not valid.")
	}
	if DataUser.Name == "" {
		return "", errors.New("Your Account not valid.")
	}
	GenOTP := util.GenerateNumber(4)
	// send generate password
	mailservice := &useemailauth.Forgot{
		Email: DataUser.Email,
		Name:  DataUser.Name,
		OTP:   GenOTP,
	}

	// check data redis
	data := redisdb.GetSession(dataForgot.Account + "_Forgot")
	if data != "" {
		redisdb.TurncateList(dataForgot.Account + "_Forgot")
	}
	//store to redis
	err = redisdb.AddSession(dataForgot.Account+"_Forgot", GenOTP, 24*time.Hour)
	if err != nil {
		return "", err
	}

	go mailservice.SendForgot()
	// err = mailservice.SendForgot()
	// if err != nil {
	// 	return "", err
	// }

	return GenOTP, nil
}

func (u *useAuht) GenOTP(ctx context.Context, dataForgot *models.ForgotForm) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	DataUser, err := u.repoAuth.GetByAccount(dataForgot.Account) //u.repoUser.GetByEmailSaUser(dataLogin.UserName)
	if err != nil {
		// return util.GoutputErrCode(http.StatusUnauthorized, "Your User/Email not valid.") //appE.ResponseError(util.GetStatusCode(err), fmt.Sprintf("%v", err), nil)
		return "", errors.New("Your Account not valid.")
	}
	if DataUser.Name == "" {
		return "", errors.New("Your Account not valid.")
	}
	GenCode := util.GenerateNumber(4)

	// send generate code
	mailService := &useemailauth.Register{
		Email:      DataUser.Email,
		Name:       DataUser.Name,
		PasswordCd: GenCode,
	}

	go mailService.SendRegister()
	// err = mailService.SendRegister()
	// if err != nil {
	// 	return output, err
	// }
	if DataUser.UserID > 0 {
		redisdb.TurncateList(dataForgot.Account + "_Register")
	}
	//store to redis
	err = redisdb.AddSession(dataForgot.Account+"_Register", GenCode, 24*time.Hour)
	if err != nil {
		return "", err
	}
	out := map[string]interface{}{
		"otp":     GenCode,
		"account": dataForgot.Account,
	}

	return out, nil
}

func (u *useAuht) ResetPassword(ctx context.Context, dataReset *models.ResetPasswd) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	if dataReset.Passwd != dataReset.ConfirmPasswd {
		return errors.New("Password dan Confirm Password harus sama.")
	}

	DataUser, err := u.repoAuth.GetByAccount(dataReset.Account)
	if err != nil {
		return err
	}

	if util.ComparePassword(DataUser.Password, util.GetPassword(dataReset.Passwd)) {
		return errors.New("Password baru tidak boleh sama dengan yang lama.")
	}

	DataUser.Password, _ = util.Hash(dataReset.Passwd)

	err = u.repoAuth.UpdatePasswordByEmail(dataReset.Account, DataUser.Password) //u.repoAuth.Update(DataUser.UserID, data)
	if err != nil {
		return err
	}

	return nil
}

func (u *useAuht) Register(ctx context.Context, dataRegister models.RegisterForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		User models.SsUser
	)

	CekData, err := u.repoAuth.GetByAccount(dataRegister.Account) //repoAuth.GetByAccount(dataRegister.EmailAddr, true)

	if CekData.Email == dataRegister.Account {
		if CekData.IsActive {
			return output, errors.New("email sudah terdaftar.")
		}
	}

	if dataRegister.Passwd != dataRegister.ConfirmPasswd {
		return output, errors.New("Password dan confirm password harus sama.")
	}
	User.Name = dataRegister.Name
	User.Password, _ = util.Hash(dataRegister.Passwd)
	User.JoinDate = time.Now()
	User.UserType = "user"
	User.IsActive = false
	User.Email = dataRegister.Account

	if CekData.UserID > 0 {
		CekData.Name = User.Name
		CekData.Password = User.Password
		CekData.JoinDate = User.JoinDate
		CekData.UserType = User.UserType
		CekData.IsActive = User.IsActive
		CekData.Email = User.Email
		err = u.repoAuth.Update(CekData.UserID, CekData)
		if err != nil {
			return output, err
		}

	} else {
		User.UserEdit = dataRegister.Name
		User.UserInput = dataRegister.Name
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
	if CekData.UserID > 0 {
		redisdb.TurncateList(dataRegister.Account + "_Register")
	}
	//store to redis
	err = redisdb.AddSession(dataRegister.Account+"_Register", GenCode, 24*time.Hour)
	if err != nil {
		return output, err
	}
	out := map[string]interface{}{
		"otp":     GenCode,
		"account": User.Email,
	}
	return out, nil
}
func (u *useAuht) VerifyRegister(ctx context.Context, dataVerify models.VerifyForm) (output interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	var (
		expireToken = setting.FileConfigSetting.JWTExpired
	)

	data := redisdb.GetSession(dataVerify.Account + "_Register")
	if data == "" {
		return output, errors.New("Account yang anda masukan salah")
	}

	if data != dataVerify.VerifyCode {
		return output, errors.New("OTP yang anda masukan salah.")
	}

	redisdb.TurncateList(dataVerify.Account + "_Register")

	DataUser, err := u.repoAuth.GetByAccount(dataVerify.Account)
	if err != nil {
		return output, errors.New("Account yang anda masukan tidak terdaftar")
	}

	DataFile, err := u.repoFile.GetBySaFileUpload(ctx, DataUser.FileID)

	token, err := util.GenerateToken(DataUser.UserID, dataVerify.Account, DataUser.UserType)
	if err != nil {
		return nil, err
	}

	mUser := map[string]interface{}{
		"is_active": true,
	}
	err = u.repoAuth.Update(DataUser.UserID, mUser)
	if err != nil {
		return output, err
	}

	redisdb.AddSession(token, DataUser.UserID, time.Duration(expireToken)*time.Hour)
	// expired FCM
	if dataVerify.FcmToken != "" {
		redisdb.AddSession(strconv.Itoa(DataUser.UserID)+"_fcm", dataVerify.FcmToken, time.Duration(expireToken)*time.Hour)
	}

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
		"user_type": "user",
	}

	return response, nil
}
func (u *useAuht) Verify(ctx context.Context, dataVerify models.VerifyForm) (err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeOut)
	defer cancel()

	data := redisdb.GetSession(dataVerify.Account + "_Forgot")
	if data == "" {
		return errors.New("Please Resend Code")
	}

	if data != dataVerify.VerifyCode {
		return errors.New("Invalid Code.")
	}
	redisdb.TurncateList(dataVerify.Account + "_Forgot")

	return nil
}
