package customer

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *CustomerMutation) DeleteCustomer(
	ctx context.Context,
	input gqlgen.DeleteCustomerInput,
) (*gqlgen.DeleteCustomerPayload, error) {
	_, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &input.ID,
	}).Company().Exec(ctx)

	if err != nil {
		return nil, err
	}

	// TODO: revert this back to soft delete after testing
	// if err := permissions.CanAccessCompany(ctx, company.ID, r.Prisma, allowedTypes); err != nil {
	// 	return nil, err
	// }

	// deleted := true

	user, err := r.Prisma.DeleteUser(prisma.UserWhereUniqueInput{
		ID: &input.ID,
	}).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: revert this back to soft delete after testing
	// user, err := r.Prisma.UpdateUser(prisma.UserUpdateParams{
	// 	Where: prisma.UserWhereUniqueInput{
	// 		ID: &input.ID,
	// 	},
	// 	Data: prisma.UserUpdateInput{
	// 		Deleted: &deleted,
	// 	},
	// }).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.DeleteCustomerPayload{
		Customer: &prisma.Customer{User: user},
	}, nil
}
