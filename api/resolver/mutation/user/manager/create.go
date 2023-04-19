package manager

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/api/resolver/mutation/auth"
	"github.com/steebchen/keskin-api/api/resolver/mutation/user/iuser"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/hours"
	"github.com/steebchen/keskin-api/lib/users"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ManagerMutation) CreateManager(
	ctx context.Context,
	input gqlgen.CreateManagerInput,
) (*gqlgen.CreateManagerPayload, error) {
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

	create := iuser.CreateUserInput(input.User, activateToken)

	activated := true

	create.Type = prisma.UserTypeManager
	create.Activated = &activated

	create.Company = &prisma.CompanyCreateOneWithoutUsersInput{
		Connect: &prisma.CompanyWhereUniqueInput{
			ID: &input.Company,
		},
	}

	user, err := r.Prisma.CreateUser(create).Exec(ctx)

	if err != nil {
		return nil, err
	}

	err = hours.UpsertWorkingHours(r.Prisma, ctx, user.ID, input.Manager.WorkingHours)

	if err != nil {
		return nil, err
	}

	return &gqlgen.CreateManagerPayload{
		Manager: &prisma.Manager{User: user},
	}, nil
}
