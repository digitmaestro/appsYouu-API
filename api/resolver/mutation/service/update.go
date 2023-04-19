package service

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/api/resolver/product_service_attribute"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/file"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ServiceMutation) UpdateService(
	ctx context.Context,
	input gqlgen.UpdateServiceInput,
	language *string,
) (*gqlgen.UpdateServicePayload, error) {
	branch, err := r.Prisma.Service(prisma.ServiceWhereUniqueInput{
		ID: &input.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	service, err := r.Prisma.Service(prisma.ServiceWhereUniqueInput{
		ID: &input.ID,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if service.Deleted {
		return nil, gqlerrors.NewPermissionError("Service is deleted")
	}

	imageID, err := file.MaybeUpload(input.Patch.Image, true)

	if err != nil {
		return nil, err
	}

	service, err = r.Prisma.UpdateService(prisma.ServiceUpdateParams{
		Where: prisma.ServiceWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.ServiceUpdateInput{
			Name:         i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Name),
			Desc:         i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Desc),
			Price:        prisma.Price(input.Patch.Price),
			Duration:     prisma.Int32Ptr(input.Patch.Duration),
			Image:        imageID,
			GenderTarget: input.Patch.GenderTarget,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	err = product_service_attribute.UpsertAttributes(r.Prisma, ctx, nil, &input.ID, input.Patch.Attributes)

	if err != nil {
		return nil, err
	}

	for _, key := range input.Patch.RemoveAttributes {
		_, err := r.Prisma.DeleteManyProductServiceAttributes(&prisma.ProductServiceAttributeWhereInput{
			Service: &prisma.ServiceWhereInput{
				ID: &input.ID,
			},
			Key: &key,
		}).Exec(ctx)

		if err != nil {
			return nil, err
		}
	}

	return &gqlgen.UpdateServicePayload{
		Service: service,
	}, nil
}
