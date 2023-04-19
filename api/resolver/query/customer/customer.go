package customer

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type CustomerQuery struct {
	Prisma *prisma.Client
}

var typeCustomer = prisma.UserTypeCustomer

var allowedTypes = []prisma.UserType{
	prisma.UserTypeEmployee,
	prisma.UserTypeManager,
}

func (r *CustomerQuery) Customer(ctx context.Context, id string, language *string) (*prisma.Customer, error) {
	company, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &id,
	}).Company().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessCompany(ctx, company.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	customer, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)

	if err != nil {
		return nil, gqlerrors.NewFormatNodeError(err, id)
	}

	return &prisma.Customer{User: customer}, nil
}

func (r *CustomerQuery) Customers(ctx context.Context, input gqlgen.CustomerInput, language *string) (*gqlgen.CustomerConnection, error) {
	deleted := false

	where := &prisma.UserWhereInput{
		Type:    &typeCustomer,
		Deleted: &deleted,
	}

	if input.Branch == nil {
		companyId := sessctx.CompanyWithFallback(ctx, r.Prisma, input.Company)

		if err := permissions.CanAccessCompany(ctx, companyId, r.Prisma, allowedTypes); err != nil {
			return nil, err
		}

		where.Company = &prisma.CompanyWhereInput{
			ID: &companyId,
		}
	} else {
		if err := permissions.CanAccessBranch(ctx, *input.Branch, r.Prisma, allowedTypes); err != nil {
			return nil, err
		}

		where.Company = &prisma.CompanyWhereInput{
			BranchesSome: &prisma.BranchWhereInput{
				ID: input.Branch,
			},
		}
	}

	users, err := r.Prisma.Users(&prisma.UsersParams{
		Where:   where,
		OrderBy: AssembleOrder(input.Order),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	customers := []*prisma.Customer{}

	for _, user := range users {
		u := user
		customers = append(customers, &prisma.Customer{User: &u})
	}

	return &gqlgen.CustomerConnection{
		Nodes: customers,
	}, nil
}

func New(client *prisma.Client) *CustomerQuery {
	return &CustomerQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
