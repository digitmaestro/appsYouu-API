package notification

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/api/resolver/mutation/email_template"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Notification) NotifyAllCustomers(
	ctx context.Context,
	input gqlgen.NotifyAllCustomersInput,
) (*gqlgen.NotificationsPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if viewer.Type != prisma.UserTypeManager && viewer.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError("user is not allowed to send notifications")
	}

	if viewer.Type == prisma.UserTypeManager && input.Company != nil {
		if err := permissions.CanAccessCompany(ctx, *input.Company, r.Prisma, []prisma.UserType{prisma.UserTypeManager}); err != nil {
			return nil, err
		}
	}

	payloads := []*gqlgen.NotificationPayload{}

	deleted := false

	companyId := sessctx.CompanyWithFallback(ctx, r.Prisma, input.Company)
	typeCustomer := prisma.UserTypeCustomer
	where := &prisma.UserWhereInput{
		Type: &typeCustomer,
		Company: &prisma.CompanyWhereInput{
			ID: &companyId,
		},
		Deleted: &deleted,
	}

	users, err := r.Prisma.Users(&prisma.UsersParams{
		Where: where,
	}).Exec(ctx)

	for _, user := range users {
		authorized, err := MaySendMessageTo(r.Prisma, ctx, viewer, &user)

		if err != nil {
			return nil, err
		}

		if authorized {
			response, _ := Send(r.Prisma, r.MessagingClient, ctx, user.ID, input.Title, input.Text)

			payloads = append(payloads, &gqlgen.NotificationPayload{
				UserID: user.ID,
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
				UserID: user.ID,
				Result: "not authorized",
			})
		}
	}

	return &gqlgen.NotificationsPayload{
		Payloads: payloads,
	}, nil
}
