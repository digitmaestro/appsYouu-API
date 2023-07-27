package validator

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
	validator "gopkg.in/go-playground/validator.v8"
)

var validate *validator.Validate = validator.New(&validator.Config{TagName: "validate"})

// type Number struct {
// 	Phone string `validate:required,number`
// }

func Phone(number string) error {
	num, err := phonenumbers.Parse(number, "PK")
	if err != nil {
		return errInvalidPhoneNumber
	}

	fmt.Println(num)

	return nil
}

func Email(email string) error {
	err := validate.Field(email, "required,email")
	if err != nil {
		return errInvalidEmail
	}

	return nil
}
