package manager

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/api/resolver/mutation/user/iuser"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/hours"
	"github.com/steebchen/keskin-api/lib/users"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ManagerMutation) UpdateManager(
	ctx context.Context,
	input gqlgen.UpdateManagerInput,
) (*gqlgen.UpdateManagerPayload, error) {
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

	if input.Patch.Email != nil {
		emailInUse, err := users.EmailInUse(ctx, r.Prisma, *input.Patch.Email, &company.ID, nil, &input.ID)

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

	err = hours.UpsertWorkingHours(r.Prisma, ctx, user.ID, input.PatchManager.WorkingHours)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateManagerPayload{
		Manager: &prisma.Manager{User: user},
	}, nil
}
