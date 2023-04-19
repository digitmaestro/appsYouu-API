package notification

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/sessctx"
)

func (r *Notification) TriggerTestNotification(
	ctx context.Context,
) (*gqlgen.NotificationPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	response, err := Send(r.Prisma, r.MessagingClient, ctx, viewer.ID, i18n.Language(ctx)["TEST_NOTIFICATION_TITLE"], i18n.Language(ctx)["TEST_NOTIFICATION_TEXT"])

	return &gqlgen.NotificationPayload{
		UserID: viewer.ID,
		Result: response,
	}, nil
}
