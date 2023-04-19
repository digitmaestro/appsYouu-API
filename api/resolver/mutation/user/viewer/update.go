package viewer

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/mutation/user/iuser"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/auth"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/lib/users"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ViewerMutation) UpdateViewer(
	ctx context.Context,
	input gqlgen.UpdateViewerInput,
) (*gqlgen.UpdateViewerPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if input.Patch.Email != nil {
		var branchId *string = nil
		var companyId *string = nil

		branch, err := r.Prisma.User(prisma.UserWhereUniqueInput{
			ID: &viewer.ID,
		}).Branch().Exec(ctx)

		if err == nil {
			branchId = &branch.ID
		}

		company, err := r.Prisma.User(prisma.UserWhereUniqueInput{
			ID: &viewer.ID,
		}).Company().Exec(ctx)

		if err == nil {
			companyId = &company.ID
		}

		emailInUse, err := users.EmailInUse(ctx, r.Prisma, *input.Patch.Email, companyId, branchId, &viewer.ID)

		if err != nil {
			return nil, err
		}

		if emailInUse {
			return nil, gqlerrors.NewValidationError("Email already used for another account", "DuplicateEmail")
		}
	}

	update := iuser.UpdateUserInput(input.Patch)
	update.AllowReviewSharing = input.AllowReviewSharing
	if input.Password != nil {
		passwordHash := auth.HashPassword(*input.Password)
		update.PasswordHash = &passwordHash
	}

	user, err := r.Prisma.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			ID: &viewer.ID,
		},
		Data: update,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateViewerPayload{
		User: user.Convert(),
	}, nil
}
