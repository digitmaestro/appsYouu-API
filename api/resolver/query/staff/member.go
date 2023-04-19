package staff

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *StaffQuery) StaffMember(ctx context.Context, id *string) (prisma.IStaff, error) {
	var allowedTypes = []prisma.UserType{
		prisma.UserTypeManager,
	}

	// if an id is provided, we have to check for manager permission
	if id != nil {
		branch, err := r.Prisma.User(prisma.UserWhereUniqueInput{
			ID: id,
		}).Branch().Exec(ctx)

		if err != nil && err != prisma.ErrNoResult {
			return nil, err
		}

		if branch != nil {
			if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, allowedTypes); err != nil {
				return nil, err
			}
		}

		company, err := r.Prisma.User(prisma.UserWhereUniqueInput{
			ID: id,
		}).Company().Exec(ctx)

		if err != nil && err != prisma.ErrNoResult {
			return nil, err
		}

		if company != nil {
			if err := permissions.CanAccessCompany(ctx, company.ID, r.Prisma, allowedTypes); err != nil {
				return nil, err
			}
		}
	}

	// if the id is nil, default to the viewer
	if id == nil {
		viewer, err := sessctx.User(ctx)

		if err != nil {
			return nil, err
		}

		id = &viewer.ID
	}

	staff, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: id,
	}).Exec(ctx)

	if err != nil {
		return nil, gqlerrors.NewFormatNodeError(err, *id)
	}

	if staff.Type != prisma.UserTypeEmployee && staff.Type != prisma.UserTypeManager {
		return nil, gqlerrors.NewInternalError("requested user is not of type employee or manager", "WrongUserType")
	}

	return staff.ConvertStaff(), nil
}
