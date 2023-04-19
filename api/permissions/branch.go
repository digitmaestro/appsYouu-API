package permissions

import (
	"context"
	"fmt"

	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func CanAccessBranch(ctx context.Context, branch string, client *prisma.Client, types []prisma.UserType) error {
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
		branches, err := client.Branches(&prisma.BranchesParams{
			Where: &prisma.BranchWhereInput{
				ID: &branch,
				EmployeesSome: &prisma.UserWhereInput{
					ID: &viewer.ID,
				},
			},
		}).Exec(ctx)

		if err != nil {
			return err
		}

		if len(branches) == 0 {
			return gqlerrors.NewPermissionError("user is not allowed to access this branch")
		}

	case prisma.UserTypeManager:
		branches, err := client.Branches(&prisma.BranchesParams{
			Where: &prisma.BranchWhereInput{
				ID: &branch,
				Company: &prisma.CompanyWhereInput{
					UsersSome: &prisma.UserWhereInput{
						ID: &viewer.ID,
					},
				},
			},
		}).Exec(ctx)

		if err != nil {
			return err
		}

		if len(branches) == 0 {
			return gqlerrors.NewPermissionError("user is not allowed to access this branch")
		}

	default:
		return gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	return nil
}
