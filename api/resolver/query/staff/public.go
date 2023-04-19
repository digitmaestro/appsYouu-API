package staff

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *StaffQuery) PublicStaffMembers(ctx context.Context, input gqlgen.StaffMemberInput) (*gqlgen.PublicStaffMemberConnection, error) {
	var allowedTypes = []prisma.UserType{
		prisma.UserTypeManager,
		prisma.UserTypeEmployee,
	}

	deleted := false

	where := &prisma.UserWhereInput{
		TypeIn:  allowedTypes,
		Deleted: &deleted,
	}

	if input.Branch == nil {
		companyId := sessctx.CompanyWithFallback(ctx, r.Prisma, input.Company)

		where.Or = []prisma.UserWhereInput{
			{
				Branch: &prisma.BranchWhereInput{
					Company: &prisma.CompanyWhereInput{
						ID: &companyId,
					},
				},
			},
			{
				Company: &prisma.CompanyWhereInput{
					ID: &companyId,
				},
			},
		}
	} else {
		where.Or = []prisma.UserWhereInput{
			{
				Branch: &prisma.BranchWhereInput{
					ID: input.Branch,
				},
			},
			{
				Company: &prisma.CompanyWhereInput{
					BranchesSome: &prisma.BranchWhereInput{
						ID: input.Branch,
					},
				},
			},
		}
	}

	users, err := r.Prisma.Users(&prisma.UsersParams{
		Where: where,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	var staff []prisma.IPublicStaff

	for _, user := range users {
		u := user
		staff = append(staff, u.ConvertPublicStaff())
	}

	return &gqlgen.PublicStaffMemberConnection{
		Nodes: staff,
	}, nil
}
