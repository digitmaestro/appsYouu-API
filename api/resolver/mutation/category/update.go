package category

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *CategoryMutation) UpdateCategory(ctx context.Context, input gqlgen.UpdateCategoryInput) (*gqlgen.UpdateCategoryPayload, error) {
	user, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if user.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	cg, err := r.Prisma.UpdateCategory(prisma.CategoryUpdateParams{
		Where: prisma.CategoryWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.CategoryUpdateInput{
			Name: input.Patch.Name,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateCategoryPayload{
		Category: cg,
	}, nil
}
