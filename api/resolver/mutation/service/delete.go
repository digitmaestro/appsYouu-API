package service

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ServiceMutation) DeleteService(
	ctx context.Context,
	input gqlgen.DeleteServiceInput,
	language *string,
) (*gqlgen.DeleteServicePayload, error) {
	branch, err := r.Prisma.Service(prisma.ServiceWhereUniqueInput{
		ID: &input.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	deleted := true

	service, err := r.Prisma.UpdateService(prisma.ServiceUpdateParams{
		Where: prisma.ServiceWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.ServiceUpdateInput{
			Deleted: &deleted,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.DeleteServicePayload{
		Service: service,
	}, nil
}
