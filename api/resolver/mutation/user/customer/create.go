package customer

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/api/resolver/mutation/auth"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/users"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *CustomerMutation) CreateCustomer(
	ctx context.Context,
	input gqlgen.CreateCustomerInput,
) (*gqlgen.CreateCustomerPayload, error) {
	if err := permissions.CanAccessCompany(ctx, input.Company, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	emailInUse, err := users.EmailInUse(ctx, r.Prisma, input.User.Email, &input.Company, nil, nil)

	if err != nil {
		return nil, err
	}

	if emailInUse {
		return nil, gqlerrors.NewValidationError("Email already used for another account", "DuplicateEmail")
	}

	activateToken, err := auth.GenerateActivateToken(r.Prisma, ctx)

	if err != nil {
		return nil, err
	}

	create := CreateUserInput(input.User, input.Customer, activateToken)

	create.Type = prisma.UserTypeCustomer

	create.Company = &prisma.CompanyCreateOneWithoutUsersInput{
		Connect: &prisma.CompanyWhereUniqueInput{
			ID: &input.Company,
		},
	}

	user, err := r.Prisma.CreateUser(create).Exec(ctx)

	if err != nil {
		return nil, err
	}

	auth.RequestPasswordReset(r.Prisma, ctx, user.Email, &input.Company)
	auth.RequestActivationLink(r.Prisma, ctx, user.Email, &input.Company)

	return &gqlgen.CreateCustomerPayload{
		Customer: &prisma.Customer{User: user},
	}, nil
}
