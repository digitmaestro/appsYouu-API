package viewer

import (
	"context"

	"github.com/99designs/gqlgen/graphql"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/file"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ViewerMutation) UploadViewerProfilePicture(ctx context.Context, upload graphql.Upload) (*gqlgen.UploadViewerProfilePicturePayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	imageID, err := file.Upload(upload, false)

	if err != nil {
		return nil, err
	}

	user, err := r.Prisma.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			ID: &viewer.ID,
		},
		Data: prisma.UserUpdateInput{
			Image: imageID,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UploadViewerProfilePicturePayload{
		Image: picture.FromID(user.Image),
		User:  user.Convert(),
	}, nil
}
