package validator

import "errors"

var (
	errInvalidPhoneNumber = errors.New("error invalid phone number")
	errInvalidEmail       = errors.New("error invalid email")
)
