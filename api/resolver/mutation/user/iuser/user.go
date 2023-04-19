package iuser

import (
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func UpdateUserInput(patch *gqlgen.UpdateUserPatch) prisma.UserUpdateInput {
	birthdate := ""

	if patch.Birthday != nil {
		birthdate = (*patch.Birthday)[5:10]
	}

	return prisma.UserUpdateInput{
		Email:       patch.Email,
		FirstName:   patch.FirstName,
		LastName:    patch.LastName,
		Gender:      patch.Gender,
		PhoneNumber: patch.PhoneNumber,
		ZipCode:     patch.ZipCode,
		Street:      patch.Street,
		City:        patch.City,
		Birthday:    patch.Birthday,
		Birthdate:   &birthdate,
	}
}

func CreateUserInput(patch *gqlgen.CreateUserData, activateToken string) prisma.UserCreateInput {
	birthdate := ""

	if patch.Birthday != nil {
		birthdate = (*patch.Birthday)[5:10]
	}

	return prisma.UserCreateInput{
		Email:         patch.Email,
		FirstName:     patch.FirstName,
		LastName:      patch.LastName,
		Gender:        &patch.Gender,
		PhoneNumber:   patch.PhoneNumber,
		ZipCode:       patch.ZipCode,
		Street:        patch.Street,
		City:          patch.City,
		Birthday:      patch.Birthday,
		Birthdate:     &birthdate,
		ActivateToken: &activateToken,
	}
}
