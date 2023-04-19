package customer

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/users"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *CustomerMutation) UpdateCustomer(
	ctx context.Context,
	input gqlgen.UpdateCustomerInput,
) (*gqlgen.UpdateCustomerPayload, error) {
	company, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &input.ID,
	}).Company().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessCompany(ctx, company.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	user, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &input.ID,
	}).Exec(ctx)

	if user.Deleted {
		return nil, gqlerrors.NewPermissionError("user is deleted")
	}

	if input.PatchUser.Email != nil {
		emailInUse, err := users.EmailInUse(ctx, r.Prisma, *input.PatchUser.Email, &company.ID, nil, &input.ID)

		if err != nil {
			return nil, err
		}

		if emailInUse {
			return nil, gqlerrors.NewValidationError("Email already used for another account", "DuplicateEmail")
		}
	}

	user, err = r.Prisma.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			ID: &input.ID,
		},
		Data: UpdateUserInput(input.PatchUser, input.PatchCustomer),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateCustomerPayload{
		Customer: &prisma.Customer{User: user},
	}, nil
}
