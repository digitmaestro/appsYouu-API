package administrator

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/mutation/auth"
	"github.com/steebchen/keskin-api/api/resolver/mutation/user/iuser"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/lib/users"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *AdministratorMutation) CreateAdministrator(
	ctx context.Context,
	input gqlgen.CreateAdministratorInput,
) (*gqlgen.CreateAdministratorPayload, error) {
	user, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if user.Deleted || user.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError("user is not allowed to create administrators")
	}

	emailInUse, err := users.EmailInUse(ctx, r.Prisma, input.User.Email, nil, nil, nil)

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

	create := iuser.CreateUserInput(input.User, activateToken)

	create.Type = prisma.UserTypeAdministrator

	user, err = r.Prisma.CreateUser(create).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.CreateAdministratorPayload{
		Administrator: &prisma.Administrator{User: user},
	}, nil
}
