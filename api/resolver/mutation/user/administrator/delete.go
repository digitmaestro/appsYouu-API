package administrator

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *AdministratorMutation) DeleteAdministrator(
	ctx context.Context,
	input gqlgen.DeleteAdministratorInput,
) (*gqlgen.DeleteAdministratorPayload, error) {
	user, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if user.Deleted || user.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError("user is not allowed to delete administrators")
	}

	deleted := true

	user, err = r.Prisma.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.UserUpdateInput{
			Deleted: &deleted,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.DeleteAdministratorPayload{
		Administrator: &prisma.Administrator{User: user},
	}, nil
}
