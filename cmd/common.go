package cmd

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

const (
	TABLE_WIDTH = 80
)

func isUsername(username string) bool {
	reg := regexp.MustCompile("^[a-zA-Z]+$")
	return reg.MatchString(username)
}

func isEmail(email string) bool {
	return govalidator.IsEmail(email)
}
