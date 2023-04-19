package manager

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ManagerMutation) DeleteManager(
	ctx context.Context,
	input gqlgen.DeleteManagerInput,
) (*gqlgen.DeleteManagerPayload, error) {
	company, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &input.ID,
	}).Company().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessCompany(ctx, company.ID, r.Prisma, allowedTypes); err != nil {
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

	return &gqlgen.DeleteManagerPayload{
		Manager: &prisma.Manager{User: user},
	}, nil
}
