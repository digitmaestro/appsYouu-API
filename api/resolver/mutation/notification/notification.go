package notification

import (
	"context"

	"firebase.google.com/go/messaging"

	"github.com/steebchen/keskin-api/prisma"
)

func MaySendMessageTo(prismaClient *prisma.Client, ctx context.Context, viewer *prisma.User, user *prisma.User) (bool, error) {
	if viewer.Type == prisma.UserTypeAdministrator {
		return true, nil
	}

	if viewer.Deleted || user.Deleted || viewer.Type != prisma.UserTypeManager {
		return false, nil
	}

	userCompany, err := prismaClient.User(prisma.UserWhereUniqueInput{
		ID: &user.ID,
	}).Company().Exec(ctx)

	if err == prisma.ErrNoResult {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	viewerCompany, err := prismaClient.User(prisma.UserWhereUniqueInput{
		ID: &viewer.ID,
	}).Company().Exec(ctx)

	if err == prisma.ErrNoResult {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return viewerCompany.ID == userCompany.ID, nil
}

func Send(prismaClient *prisma.Client, messagingClient *messaging.Client, ctx context.Context, userId string, title string, body string) (string, error) {
	user, err := prismaClient.User(prisma.UserWhereUniqueInput{
		ID: &userId,
	}).Exec(ctx)

	if err != nil {
		return "", err
	}

	if user.Deleted {
		return "Cannot send notifications to deleted users", nil
	}

	if user.NotificationToken != nil {
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: title,
				Body:  body,
			},
			Token: *user.NotificationToken,
		}

		response, err := messagingClient.Send(ctx, message)

		if err != nil {
			response = err.Error()
		}

		return response, err
	} else {
		return "invalid notification token", nil
	}
}

type Notification struct {
	Prisma          *prisma.Client
	MessagingClient *messaging.Client
}

func New(client *prisma.Client, messagingClient *messaging.Client) *Notification {
	return &Notification{
		Prisma:          client,
		MessagingClient: messagingClient,
	}
}
