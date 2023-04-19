package staff

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *StaffQuery) StaffMembers(ctx context.Context, input gqlgen.StaffMemberInput) (*gqlgen.StaffMemberConnection, error) {
	var allowedTypes = []prisma.UserType{
		prisma.UserTypeManager,
		prisma.UserTypeEmployee,
	}

	var users = []prisma.User{}
	var err error

	company := input.Company

	if company == nil && input.Branch == nil {
		userCompany := sessctx.Company(ctx)
		company = &userCompany
	}

	deleted := false

	if company != nil {
		if err = permissions.CanAccessCompany(ctx, *company, r.Prisma, allowedTypes); err != nil {
			return nil, err
		}

		users, err = r.Prisma.Users(&prisma.UsersParams{
			Where: &prisma.UserWhereInput{
				TypeIn:  allowedTypes,
				Deleted: &deleted,
				Or: []prisma.UserWhereInput{
					{
						Branch: &prisma.BranchWhereInput{
							Company: &prisma.CompanyWhereInput{
								ID: company,
							},
						},
					},
					{
						Company: &prisma.CompanyWhereInput{
							ID: company,
						},
					},
				},
			},
		}).Exec(ctx)
	}

	if input.Branch != nil {
		if err = permissions.CanAccessBranch(ctx, *input.Branch, r.Prisma, allowedTypes); err != nil {
			return nil, err
		}

		users, err = r.Prisma.Users(&prisma.UsersParams{
			Where: &prisma.UserWhereInput{
				TypeIn:  allowedTypes,
				Deleted: &deleted,
				Or: []prisma.UserWhereInput{
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
				},
			},
		}).Exec(ctx)
	}

	if err != nil {
		return nil, err
	}

	var staff []prisma.IStaff

	for _, user := range users {
		u := user
		staff = append(staff, u.ConvertStaff())
	}

	return &gqlgen.StaffMemberConnection{
		Nodes: staff,
	}, nil
}
