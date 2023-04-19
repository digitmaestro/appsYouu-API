package product

import (
	"context"
	"fmt"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ProductMutation) DeleteProduct(
	ctx context.Context,
	input gqlgen.DeleteProductInput,
	language *string,
) (*gqlgen.DeleteProductPayload, error) {
	fmt.Println(input.ID)
	branch, err := r.Prisma.Product(prisma.ProductWhereUniqueInput{
		ID: &input.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	deleted := true

	product, err := r.Prisma.UpdateProduct(prisma.ProductUpdateParams{
		Where: prisma.ProductWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.ProductUpdateInput{
			Deleted: &deleted,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.DeleteProductPayload{
		Product: product,
	}, nil
}
