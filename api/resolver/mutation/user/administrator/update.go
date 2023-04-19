package administrator

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/mutation/user/iuser"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/lib/users"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *AdministratorMutation) UpdateAdministrator(
	ctx context.Context,
	input gqlgen.UpdateAdministratorInput,
) (*gqlgen.UpdateAdministratorPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if viewer.Deleted || viewer.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError("user is not allowed to modify administrators")
	}

	user, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &input.ID,
	}).Exec(ctx)

	if user.Deleted {
		return nil, gqlerrors.NewPermissionError("user is deleted")
	}

	if input.Patch.Email != nil {
		emailInUse, err := users.EmailInUse(ctx, r.Prisma, *input.Patch.Email, nil, nil, &input.ID)

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
		Data: iuser.UpdateUserInput(input.Patch),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateAdministratorPayload{
		Administrator: &prisma.Administrator{User: user},
	}, nil
}
