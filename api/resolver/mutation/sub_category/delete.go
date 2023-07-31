package sub_category

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *SubCategoryMutation) DeleteSubCategory(ctx context.Context, input gqlgen.DeleteSubCategoryInput) (*gqlgen.DeleteSubCategoryPayload, error) {
	user, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if user.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	subcg, err := r.Prisma.DeleteSubCategory(prisma.SubCategoryWhereUniqueInput{
		ID: &input.ID,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.DeleteSubCategoryPayload{
		SubCategory: subcg,
	}, nil
}
