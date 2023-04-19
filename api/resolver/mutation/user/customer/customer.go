package customer

import (
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

type CustomerMutation struct {
	Prisma *prisma.Client
}

var allowedTypes = []prisma.UserType{
	prisma.UserTypeEmployee,
	prisma.UserTypeManager,
}

func New(client *prisma.Client) *CustomerMutation {
	return &CustomerMutation{
		Prisma: client,
	}
}

func UpdateUserInput(patchUser *gqlgen.UpdateUserPatch, patchCustomer *gqlgen.UpdateCustomerPatch) prisma.UserUpdateInput {
	birthdate := ""

	if patchUser.Birthday != nil {
		birthdate = (*patchUser.Birthday)[5:10]
	}

	return prisma.UserUpdateInput{
		Email:              patchUser.Email,
		FirstName:          patchUser.FirstName,
		LastName:           patchUser.LastName,
		Gender:             patchUser.Gender,
		PhoneNumber:        patchUser.PhoneNumber,
		Note:               patchCustomer.Note,
		ZipCode:            patchUser.ZipCode,
		Street:             patchUser.Street,
		City:               patchUser.City,
		Birthday:           patchUser.Birthday,
		Birthdate:          &birthdate,
		AllowReviewSharing: patchCustomer.AllowReviewSharing,
	}
}

func CreateUserInput(user *gqlgen.CreateUserData, customer *gqlgen.CreateCustomerData, activateToken string) prisma.UserCreateInput {
	birthdate := ""

	if user.Birthday != nil {
		birthdate = (*user.Birthday)[5:10]
	}

	return prisma.UserCreateInput{
		Email:              user.Email,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Gender:             &user.Gender,
		PhoneNumber:        user.PhoneNumber,
		Note:               customer.Note,
		ZipCode:            user.ZipCode,
		Street:             user.Street,
		City:               user.City,
		Birthday:           user.Birthday,
		Birthdate:          &birthdate,
		AllowReviewSharing: customer.AllowReviewSharing,
		ActivateToken:      &activateToken,
	}
}
