package useemailauth

import (
	"fmt"
	templateemail "nuryanto2121/cukur_in_user/pkg/email"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strings"
)

type Register struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	PasswordCd string `json:"password_cd"`
}

func (R *Register) SendRegister() error {
	subjectEmail := "Informasi OTP"
	fmt.Printf(subjectEmail)
	err := util.SendEmail(R.Email, subjectEmail, getVerifyBody(R))
	if err != nil {
		return err
	}
	return nil
}

func getVerifyBody(R *Register) string {
	registerHTML := templateemail.SendRegister

	registerHTML = strings.ReplaceAll(registerHTML, `{Name}`, R.Name)
	registerHTML = strings.ReplaceAll(registerHTML, `{Email}`, R.Email)
	registerHTML = strings.ReplaceAll(registerHTML, `{PasswordCode}`, R.PasswordCd)
	return registerHTML
}
