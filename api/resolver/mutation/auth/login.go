package auth

import (
	"context"
	"fmt"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/auth"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/lib/validator"
	"github.com/steebchen/keskin-api/prisma"
)

const noCompanyFoundCode = "NoCompanyFound"

func (a *Auth) Login(ctx context.Context, input gqlgen.LoginInput) (*gqlgen.LoginPayload, error) {
	companyID := sessctx.CompanyWithFallback(ctx, a.Prisma, input.Company)

	err := validator.Email(input.Email)
	if err != nil {
		return nil, gqlerrors.NewValidationError(err.Error(), "InvalidEmail")
	}

	deleted := false
	// activated := true

	where := &prisma.UserWhereInput{
		Email:   &input.Email,
		Deleted: &deleted,
		// Activated: &activated,
	}

	// if a companyID is provided, ask for it
	// if not, ignore it because it's only required by users with type Customer, Manager, Employee
	if companyID != "" {
		// if a company ID is given, we want to make sure that it's valid
		_, err := a.Prisma.Company(prisma.CompanyWhereUniqueInput{
			ID: &companyID,
		}).Exec(ctx)

		if err == prisma.ErrNoResult {
			return nil, gqlerrors.NewInternalError("no company header found", noCompanyFoundCode)
		}

		userTypeCustomer := prisma.UserTypeCustomer
		userTypeManager := prisma.UserTypeManager
		userTypeEmployee := prisma.UserTypeEmployee
		userTypeAdministrator := prisma.UserTypeAdministrator

		where.Or = []prisma.UserWhereInput{{
			Type: &userTypeEmployee,
			Branch: &prisma.BranchWhereInput{
				Company: &prisma.CompanyWhereInput{
					ID: &companyID,
				},
			},
		}, {
			TypeIn: []prisma.UserType{userTypeCustomer, userTypeManager},
			Company: &prisma.CompanyWhereInput{
				ID: &companyID,
			},
		}, {
			Type: &userTypeAdministrator,
		}}
	}

	users, err := a.Prisma.Users(&prisma.UsersParams{
		Where: where,
	}).Exec(ctx)

	if err != nil {
		fmt.Println("Error: " + err.Error())

		return nil, err
	}

	// TODO waiting for https://github.com/prisma/prisma/issues/4141
	if len(users) == 0 {
		return nil, gqlerrors.NewValidationError("No user found", "UserNotFound")
	}

	user := users[0]

	if !user.Activated {
		return nil, gqlerrors.NewVerificationError("user with email("+user.Email+")"+" and id("+user.ID+") not verified", "UserNotVerified")
	}

	// we want to make sure that users only log in on pages where they're supposed to.
	if user.Type != prisma.UserTypeAdministrator && companyID == "" {
		fmt.Println("error fkk")
		return nil, gqlerrors.NewInternalError("company id required", noCompanyFoundCode)
	}

	err = auth.VerifyPassword(user.PasswordHash, input.Password)
	if err != nil {
		return nil, gqlerrors.NewValidationError("password is incorrect", "PasswordIncorrect")
	}

	create := prisma.SessionCreateInput{
		User: prisma.UserCreateOneWithoutSessionsInput{
			Connect: &prisma.UserWhereUniqueInput{
				ID: &user.ID,
			},
		},
		Token: auth.GenerateToken(),
	}

	if companyID != "" {
		create.Company = &prisma.CompanyCreateOneInput{
			Connect: &prisma.CompanyWhereUniqueInput{
				ID: &companyID,
			},
		}
	}

	session, err := a.Prisma.CreateSession(create).Exec(ctx)

	if err != nil {
		return nil, err
	}

	loginFix := sessctx.LoginFix(ctx)
	loginFix.Apply = true
	loginFix.SessionToken = session.Token
	loginFix.User = &user

	return &gqlgen.LoginPayload{
		Session: session,
		User:    user.Convert(),
	}, nil
}
