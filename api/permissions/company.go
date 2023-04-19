package permissions

import (
	"context"
	"fmt"

	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func CanAccessCompany(ctx context.Context, company string, client *prisma.Client, types []prisma.UserType) error {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return err
	}

	if !in(viewer.Type, types) && viewer.Type != prisma.UserTypeAdministrator {
		return gqlerrors.NewPermissionError(fmt.Sprintf("user has type %s, but only types %+v are allowed", viewer.Type, types))
	}

	switch viewer.Type {
	case prisma.UserTypeAdministrator:
		return nil
	case prisma.UserTypeEmployee:
		companies, err := client.Companies(&prisma.CompaniesParams{
			Where: &prisma.CompanyWhereInput{
				ID: &company,
				BranchesSome: &prisma.BranchWhereInput{
					EmployeesSome: &prisma.UserWhereInput{
						ID: &viewer.ID,
					},
				},
			},
		}).Exec(ctx)

		if err != nil {
			return err
		}

		if len(companies) == 0 {
			return gqlerrors.NewPermissionError("user is not allowed to access this company")
		}

	case prisma.UserTypeManager:
		companies, err := client.Companies(&prisma.CompaniesParams{
			Where: &prisma.CompanyWhereInput{
				ID: &company,
				UsersSome: &prisma.UserWhereInput{
					ID: &viewer.ID,
				},
			},
		}).Exec(ctx)

		if err != nil {
			return err
		}

		if len(companies) == 0 {
			return gqlerrors.NewPermissionError("user is not allowed to access this company")
		}

	default:
		return gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	return nil
}
