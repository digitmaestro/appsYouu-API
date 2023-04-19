package company

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *CompanyMutation) DeleteCompany(
	ctx context.Context,
	input gqlgen.DeleteCompanyInput,
	language *string,
) (*gqlgen.DeleteCompanyPayload, error) {
	user, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if user.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	company, err := r.Prisma.DeleteCompany(prisma.CompanyWhereUniqueInput{
		ID: &input.ID,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.DeleteCompanyPayload{
		Company: company,
	}, nil
}
