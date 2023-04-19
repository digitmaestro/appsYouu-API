package notification

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/mutation/email_template"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Notification) SendUserNotification(
	ctx context.Context,
	input *gqlgen.SendUserNotificationInput,
) (*gqlgen.NotificationsPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if viewer.Type != prisma.UserTypeManager && viewer.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError("user is not allowed to send notifications")
	}

	payloads := []*gqlgen.NotificationPayload{}

	for _, userId := range input.Users {
		user, err := r.Prisma.User(prisma.UserWhereUniqueInput{
			ID: &userId,
		}).Exec(ctx)

		if err != nil {
			return nil, err
		}

		authorized, err := MaySendMessageTo(r.Prisma, ctx, viewer, user)

		if err != nil {
			return nil, err
		}

		if authorized {
			response, _ := Send(r.Prisma, r.MessagingClient, ctx, userId, input.Title, input.Text)

			payloads = append(payloads, &gqlgen.NotificationPayload{
				UserID: userId,
				Result: response,
			})

			branches, err := r.Prisma.Branches(&prisma.BranchesParams{
				Where: &prisma.BranchWhereInput{
					Company: &prisma.CompanyWhereInput{
						UsersSome: &prisma.UserWhereInput{
							ID: &user.ID,
						},
					},
				},
			}).Exec(ctx)

			if err == nil && len(branches) > 0 {
				go email_template.SendEmail(
					context.Background(),
					r.Prisma,
					branches[0].ID,
					user.Email,
					input.Title,
					input.Text,
				)
			}
		} else {
			payloads = append(payloads, &gqlgen.NotificationPayload{
				UserID: userId,
				Result: "not authorized",
			})
		}
	}

	return &gqlgen.NotificationsPayload{
		Payloads: payloads,
	}, nil
}
