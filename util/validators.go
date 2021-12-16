package util

import (
	"regexp"
	"github.com/gibrangul95/go-todos/internal/model"
	valid "github.com/asaskevich/govalidator"
)

func IsEmpty(str string) (bool, string) {
	if valid.HasWhitespaceOnly(str) && str != "" {
		return true, "Must not be empty"
	}

	return false, ""
}

func ValidateRegister(u * model.User) *model.UserErrors {
	e := &model.UserErrors{}
	e.Err, e.Username = IsEmpty(u.Username)

	if !valid.IsEmail(u.Email){
		e.Err, e.Email = true, "Must be a valid email"
	}

	re := regexp.MustCompile("\\d")
	if !(len(u.Password) >= 8 && valid.HasLowerCase(u.Password) && valid.HasUpperCase(u.Password) && re.MatchString(u.Password)) {
		e.Err, e.Password = true, "Length of password should be atleast 8 and it must be a combination of uppercase letters, lowercase letters and numbers"
	}

	return e
}