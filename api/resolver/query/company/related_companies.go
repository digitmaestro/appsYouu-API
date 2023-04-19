package company

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/auth"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *CompanyQuery) RelatedCompanies(ctx context.Context, email string, password string, language *string) (*gqlgen.RelatedCompanies, error) {
	deleted := false

	users, err := r.Prisma.Users(&prisma.UsersParams{
		Where: &prisma.UserWhereInput{
			Email:  &email,
			TypeIn: []prisma.UserType{prisma.UserTypeManager, prisma.UserTypeAdministrator},
			Deleted: &deleted,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// TODO waiting for https://github.com/prisma/prisma/issues/4141
	if len(users) == 0 {
		return nil, gqlerrors.NewValidationError("No user found", "UserNotFound")
	}

	validPasswordFound := false

	for _, user := range users {
		err = auth.VerifyPassword(user.PasswordHash, password)

		if err == nil {
			validPasswordFound = true
		}
	}

	if !validPasswordFound {
		return nil, gqlerrors.NewValidationError("password is incorrect", "PasswordIncorrect")
	}

	isAdmin := false
	companies := []*prisma.Company{}

	for _, user := range users {
		if user.Type == prisma.UserTypeAdministrator {
			isAdmin = true
		}

		company, err := r.Prisma.User(prisma.UserWhereUniqueInput{
			ID: &user.ID,
		}).Company().Exec(ctx)

		if err == nil && company != nil {
			clone := company
			companies = append(companies, clone)
		}
	}

	return &gqlgen.RelatedCompanies{
		Companies: companies,
		IsAdmin:   isAdmin,
	}, nil
}
