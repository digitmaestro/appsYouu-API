package product

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

func (r *ProductMutation) UpdateProduct(
	ctx context.Context,
	input gqlgen.UpdateProductInput,
	language *string,
) (*gqlgen.UpdateProductPayload, error) {
	branch, err := r.Prisma.Product(prisma.ProductWhereUniqueInput{
		ID: &input.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	product, err := r.Prisma.Product(prisma.ProductWhereUniqueInput{
		ID: &input.ID,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if product.Deleted {
		return nil, gqlerrors.NewPermissionError("Product is deleted")
	}

	imageID, err := file.MaybeUpload(input.Patch.Image, true)

	if err != nil {
		return nil, err
	}

	product, err = r.Prisma.UpdateProduct(prisma.ProductUpdateParams{
		Where: prisma.ProductWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.ProductUpdateInput{
			Name:     i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Name),
			Desc:     i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Desc),
			Price:    prisma.Price(input.Patch.Price),
			Capacity: input.Patch.Capacity,
			Category: input.Patch.Category,
			Image:    imageID,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	err = product_service_attribute.UpsertAttributes(r.Prisma, ctx, &input.ID, nil, input.Patch.Attributes)

	if err != nil {
		return nil, err
	}

	for _, key := range input.Patch.RemoveAttributes {
		_, err := r.Prisma.DeleteManyProductServiceAttributes(&prisma.ProductServiceAttributeWhereInput{
			Product: &prisma.ProductWhereInput{
				ID: &input.ID,
			},
			Key: &key,
		}).Exec(ctx)

		if err != nil {
			return nil, err
		}
	}

	return &gqlgen.UpdateProductPayload{
		Product: product,
	}, nil
}
