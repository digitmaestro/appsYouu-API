package category

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *CategoryMutation) CreateCategory(ctx context.Context, input gqlgen.CreateCategoryInput) (*gqlgen.CreateCategoryPayload, error) {
	// user, err := sessctx.User(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// if user.Type != prisma.UserTypeAdministrator {
	// 	return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	// }

	cg, err := r.Prisma.CreateCategory(prisma.CategoryCreateInput{
		Name: *input.Data.Name,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.CreateCategoryPayload{
		Category: cg,
	}, nil
}
