package viewer

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ViewerMutation) UpdateNotificationToken(
	ctx context.Context,
	input gqlgen.UpdateNotificationTokenInput,
) (*gqlgen.UpdateViewerPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	user, err := r.Prisma.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			ID: &viewer.ID,
		},
		Data: prisma.UserUpdateInput{
			NotificationToken: &input.Token,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateViewerPayload{
		User: user.Convert(),
	}, nil
}
