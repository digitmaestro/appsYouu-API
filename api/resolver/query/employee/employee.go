package employee

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type EmployeeQuery struct {
	Prisma *prisma.Client
}

var typeEmployee = prisma.UserTypeEmployee

var allowedTypes = []prisma.UserType{
	prisma.UserTypeManager,
}

func (r *EmployeeQuery) Employee(ctx context.Context, id string, language *string) (*prisma.Employee, error) {
	branch, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &id,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	employee, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)

	if err != nil {
		return nil, gqlerrors.NewFormatNodeError(err, id)
	}

	return &prisma.Employee{User: employee}, nil
}

func (r *EmployeeQuery) Employees(ctx context.Context, input gqlgen.EmployeeInput, language *string) (*gqlgen.EmployeeConnection, error) {
	deleted := false
	where := &prisma.UserWhereInput{
		Type:    &typeEmployee,
		Deleted: &deleted,
	}

	if input.Branch == nil {
		companyId := sessctx.CompanyWithFallback(ctx, r.Prisma, input.Company)

		if err := permissions.CanAccessCompany(ctx, companyId, r.Prisma, allowedTypes); err != nil {
			return nil, err
		}

		where.Branch = &prisma.BranchWhereInput{
			Company: &prisma.CompanyWhereInput{
				ID: &companyId,
			},
		}
	} else {
		if err := permissions.CanAccessBranch(ctx, *input.Branch, r.Prisma, allowedTypes); err != nil {
			return nil, err
		}

		where.Branch = &prisma.BranchWhereInput{
			ID: input.Branch,
		}
	}

	users, err := r.Prisma.Users(&prisma.UsersParams{
		Where: where,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	employees := []*prisma.Employee{}

	for _, user := range users {
		u := user
		employees = append(employees, &prisma.Employee{User: &u})
	}

	return &gqlgen.EmployeeConnection{
		Nodes: employees,
	}, nil
}

func New(client *prisma.Client) *EmployeeQuery {
	return &EmployeeQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
