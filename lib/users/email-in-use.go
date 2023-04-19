package users

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func EmailInUse(ctx context.Context, prismaClient *prisma.Client, email string, company *string, branch *string, ignoreUser *string) (bool, error) {
	administrator := prisma.UserTypeAdministrator
	manager := prisma.UserTypeManager
	employee := prisma.UserTypeEmployee
	customer := prisma.UserTypeCustomer

	companyId := company

	if companyId == nil && branch != nil {
		company, err := prismaClient.Branch(prisma.BranchWhereUniqueInput{
			ID: branch,
		}).Company().Exec(ctx)

		if err != nil {
			return true, err
		}

		companyId = &company.ID
	}

	deleted := false

	where := prisma.UserWhereInput{
		Email:   &email,
		Deleted: &deleted,
	}

	if ignoreUser != nil {
		where.IDNot = ignoreUser
	}

	if companyId != nil {
		where.Or = []prisma.UserWhereInput{
			{
				Type: &administrator,
			},
			{
				Type: &employee,
				Branch: &prisma.BranchWhereInput{
					Company: &prisma.CompanyWhereInput{
						ID: companyId,
					},
				},
			},
			{
				TypeIn: []prisma.UserType{manager, customer},
				Company: &prisma.CompanyWhereInput{
					ID: companyId,
				},
			},
		}
	}

	duplicateUsers, err := prismaClient.Users(&prisma.UsersParams{
		Where: &where,
	}).Exec(ctx)

	if err != nil {
		return true, err
	}

	if len(duplicateUsers) > 0 {
		return true, nil
	}

	return false, nil
}
