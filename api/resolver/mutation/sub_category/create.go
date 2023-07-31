package sub_category

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *SubCategoryMutation) CreateSubCategory(ctx context.Context, input gqlgen.CreateSubCategoryInput) (*gqlgen.CreateSubCategoryPayload, error) {
	cg, err := r.Prisma.Category(prisma.CategoryWhereUniqueInput{
		ID: &input.Data.CategoryID,
	}).Exec(ctx)

	if err != nil {
		return nil, gqlerrors.NewNotFoundError(input.Data.CategoryID)
	}

	subCg, err := r.Prisma.CreateSubCategory(prisma.SubCategoryCreateInput{
		Name: *input.Data.Name,
		Category: prisma.CategoryCreateOneWithoutSubCategoriesInput{
			Connect: &prisma.CategoryWhereUniqueInput{
				ID: &cg.ID,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.CreateSubCategoryPayload{
		SubCategory: subCg,
	}, nil
}
