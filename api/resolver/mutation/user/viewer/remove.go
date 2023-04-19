package viewer

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ViewerMutation) RemoveViewerProfilePicture(ctx context.Context) (*gqlgen.RemoveViewerProfilePicturePayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	user, err := r.Prisma.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			ID: &viewer.ID,
		},
		Data: prisma.UserUpdateInput{
			// temporary workaround for sql null
			Image: prisma.Str(""),
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.RemoveViewerProfilePicturePayload{
		Image: nil,
		User:  user.Convert(),
	}, nil
}
