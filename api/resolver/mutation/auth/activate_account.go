package auth

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/prisma"
)

func (a *Auth) ActivateAccount(
	ctx context.Context,
	input gqlgen.ActivateAccountInput,
) (*gqlgen.ActivateAccountPayload, error) {
	users, err := a.Prisma.Users(&prisma.UsersParams{
		Where: &prisma.UserWhereInput{
			ActivateToken: &input.Token,
		},
	}).Exec(ctx)

	if err != nil || len(users) == 0 {
		return nil, gqlerrors.NewPermissionError("invalid account activation token")
	}

	activated := true

	_, err = a.Prisma.UpdateManyUsers(prisma.UserUpdateManyParams{
		Where: &prisma.UserWhereInput{
			ActivateToken: &input.Token,
		},
		Data: prisma.UserUpdateManyMutationInput{
			Activated: &activated,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.ActivateAccountPayload{
		Status: "OK",
	}, nil
}
