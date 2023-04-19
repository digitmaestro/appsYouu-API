package employee

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

func (r *EmployeeMutation) CreateEmployee(
	ctx context.Context,
	input gqlgen.CreateEmployeeInput,
) (*gqlgen.CreateEmployeePayload, error) {
	if err := permissions.CanAccessBranch(ctx, input.Branch, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	emailInUse, err := users.EmailInUse(ctx, r.Prisma, input.Data.Email, nil, &input.Branch, nil)

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

	create := iuser.CreateUserInput(input.Data, activateToken)

	activated := true

	create.Type = prisma.UserTypeEmployee
	create.Activated = &activated

	create.Branch = &prisma.BranchCreateOneWithoutEmployeesInput{
		Connect: &prisma.BranchWhereUniqueInput{
			ID: &input.Branch,
		},
	}

	user, err := r.Prisma.CreateUser(create).Exec(ctx)

	if err != nil {
		return nil, err
	}

	err = hours.UpsertWorkingHours(r.Prisma, ctx, user.ID, input.Employee.WorkingHours)

	if err != nil {
		return nil, err
	}

	company, err := r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &input.Branch,
	}).Company().Exec(ctx)

	if err != nil {
		return nil, err
	}

	auth.RequestPasswordReset(r.Prisma, ctx, user.Email, &company.ID)

	return &gqlgen.CreateEmployeePayload{
		Employee: &prisma.Employee{User: user},
	}, nil
}
