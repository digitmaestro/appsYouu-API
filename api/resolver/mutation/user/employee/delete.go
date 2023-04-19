package employee

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *EmployeeMutation) DeleteEmployee(
	ctx context.Context,
	input gqlgen.DeleteEmployeeInput,
) (*gqlgen.DeleteEmployeePayload, error) {
	branch, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &input.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	deleted := true

	user, err := r.Prisma.UpdateUser(prisma.UserUpdateParams{
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

	return &gqlgen.DeleteEmployeePayload{
		Employee: &prisma.Employee{User: user},
	}, nil
}
