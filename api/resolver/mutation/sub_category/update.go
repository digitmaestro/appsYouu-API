package sub_category

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *SubCategoryMutation) UpdateSubCategory(ctx context.Context, input gqlgen.UpdateSubCategoryInput) (*gqlgen.UpdateSubCategoryPayload, error) {
	user, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if user.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	subCg, err := r.Prisma.UpdateSubCategory(prisma.SubCategoryUpdateParams{
		Where: prisma.SubCategoryWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.SubCategoryUpdateInput{
			Name: input.Patch.Name,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateSubCategoryPayload{
		SubCategory: subCg,
	}, nil
}
