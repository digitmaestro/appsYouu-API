package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/messaging"
	cron "github.com/robfig/cron/v3"

	"github.com/steebchen/keskin-api/api/resolver/mutation/email_template"
	"github.com/steebchen/keskin-api/api/resolver/mutation/notification"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func beginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func notifyTomorrowsClients(client *prisma.Client, messagingClient *messaging.Client) {
	ctx := context.Background()
	today := time.Now()
	day := time.Duration(24) * time.Hour

	startOfTomorrow := beginningOfDay(today.Add(day)).Format(time.RFC3339)
	endOfTomorrow := beginningOfDay(today.Add(day * 2)).Format(time.RFC3339)
	notNotified := false
	approved := prisma.AppointmentStatusApproved

	appointments, err := client.Appointments(&prisma.AppointmentsParams{
		Where: &prisma.AppointmentWhereInput{
			StartGte:         &startOfTomorrow,
			StartLt:          &endOfTomorrow,
			CustomerNotified: &notNotified,
			Status:           &approved,
		},
	}).Exec(ctx)

	if err != nil {
		return
	}

	for _, appointment := range appointments {
		customer, err := client.Appointment(prisma.AppointmentWhereUniqueInput{
			ID: &appointment.ID,
		}).Customer().Exec(ctx)

		branch, err := client.Appointment(prisma.AppointmentWhereUniqueInput{
			ID: &appointment.ID,
		}).Branch().Exec(ctx)

		ctx = sessctx.SetLanguage(ctx, customer.Language)

		if err == nil && !customer.Deleted {
			appointmentDate := i18n.FormatDate(ctx, appointment.Start)
			appointmentTime := i18n.FormatTime(ctx, appointment.Start)

			_, err := notification.Send(
				client,
				messagingClient,
				ctx,
				customer.ID,
				i18n.Language(ctx)["APPOINTMENT_REMINDER_TITLE"],
				fmt.Sprintf(
					i18n.Language(ctx)["APPOINTMENT_REMINDER_TEXT"],
					appointmentDate,
					appointmentTime,
				),
			)

			if err != nil {
				log.Println(err)
			}

			notified := true

			client.UpdateAppointment(prisma.AppointmentUpdateParams{
				Where: prisma.AppointmentWhereUniqueInput{
					ID: &appointment.ID,
				},
				Data: prisma.AppointmentUpdateInput{
					CustomerNotified: &notified,
				},
			}).Exec(ctx)

			go email_template.SendEmailTemplate(
				ctx,
				client,
				"appointmentReminder",
				branch.ID,
				customer.Email,
				customer.Gender,
				customer.LastName,
				customer.FirstName,
				&appointmentDate,
				&appointmentTime,
				nil,
				nil,
			)
		}
	}
}

func notifyOneHourBeforeAppointment(client *prisma.Client, messagingClient *messaging.Client) {
	ctx := context.Background()
	now := time.Now()
	oneHour := time.Duration(1) * time.Hour

	startOfPeriod := now.Format(time.RFC3339)
	endOfPeriod := now.Add(oneHour).Format(time.RFC3339)
	notNotified := false
	approved := prisma.AppointmentStatusApproved

	appointments, err := client.Appointments(&prisma.AppointmentsParams{
		Where: &prisma.AppointmentWhereInput{
			StartGte:                     &startOfPeriod,
			StartLt:                      &endOfPeriod,
			CustomerNotifiedAnHourBefore: &notNotified,
			Status:                       &approved,
		},
	}).Exec(ctx)

	if err != nil {
		return
	}

	for _, appointment := range appointments {
		customer, err := client.Appointment(prisma.AppointmentWhereUniqueInput{
			ID: &appointment.ID,
		}).Customer().Exec(ctx)

		branch, err := client.Appointment(prisma.AppointmentWhereUniqueInput{
			ID: &appointment.ID,
		}).Branch().Exec(ctx)

		ctx = sessctx.SetLanguage(ctx, customer.Language)

		if err == nil && !customer.Deleted {
			appointmentDate := i18n.FormatDate(ctx, appointment.Start)
			appointmentTime := i18n.FormatTime(ctx, appointment.Start)

			_, err := notification.Send(
				client,
				messagingClient,
				ctx,
				customer.ID,
				i18n.Language(ctx)["APPOINTMENT_REMINDER_TITLE"],
				fmt.Sprintf(
					i18n.Language(ctx)["APPOINTMENT_REMINDER_TEXT"],
					appointmentDate,
					appointmentTime,
				),
			)

			if err != nil {
				log.Println(err)
			}

			notified := true

			client.UpdateAppointment(prisma.AppointmentUpdateParams{
				Where: prisma.AppointmentWhereUniqueInput{
					ID: &appointment.ID,
				},
				Data: prisma.AppointmentUpdateInput{
					CustomerNotifiedAnHourBefore: &notified,
				},
			}).Exec(ctx)

			go email_template.SendEmailTemplate(
				ctx,
				client,
				"appointmentReminder",
				branch.ID,
				customer.Email,
				customer.Gender,
				customer.LastName,
				customer.FirstName,
				&appointmentDate,
				&appointmentTime,
				nil,
				nil,
			)
		}
	}
}

func notifyBirthdayClients(client *prisma.Client, messagingClient *messaging.Client) {
	ctx := context.Background()

	today := beginningOfDay(time.Now()).Format(time.RFC3339)[5:10]
	userTypeCustomer := prisma.UserTypeCustomer
	deleted := false

	users, err := client.Users(&prisma.UsersParams{
		Where: &prisma.UserWhereInput{
			Birthdate: &today,
			Type:      &userTypeCustomer,
			Deleted:   &deleted,
		},
	}).Exec(ctx)

	if err != nil {
		return
	}

	for _, user := range users {
		notificationContext := context.Background()
		notificationContext = sessctx.SetLanguage(notificationContext, user.Language)

		go notification.Send(
			client,
			messagingClient,
			notificationContext,
			user.ID,
			i18n.Language(notificationContext)["BIRTHDAY_TITLE"],
			i18n.Language(notificationContext)["BIRTHDAY_TEXT"],
		)

		company, err := client.User(prisma.UserWhereUniqueInput{
			ID: &user.ID,
		}).Company().Exec(ctx)

		var branchesParams *prisma.BranchesParams = nil

		if err == nil && company != nil {
			branchesParams = &prisma.BranchesParams{
				Where: &prisma.BranchWhereInput{
					Company: &prisma.CompanyWhereInput{
						ID: &company.ID,
					},
				},
			}
		}

		branches, err := client.Branches(branchesParams).Exec(ctx)

		if err == nil && len(branches) > 0 {
			branch := branches[0]

			go email_template.SendEmailTemplate(
				notificationContext,
				client,
				"birthday",
				branch.ID,
				user.Email,
				user.Gender,
				user.LastName,
				user.FirstName,
				nil,
				nil,
				nil,
				nil,
			)
		}
	}
}

func NewCronJobs(client *prisma.Client, messagingClient *messaging.Client) *cron.Cron {
	c := cron.New()
	c.AddFunc("CRON_TZ=Europe/Berlin 0 12 * * *", func() {
		notifyTomorrowsClients(client, messagingClient)
	})
	c.AddFunc("CRON_TZ=Europe/Berlin 0 12 * * *", func() {
		notifyBirthdayClients(client, messagingClient)
	})
	c.AddFunc("CRON_TZ=Europe/Berlin */5 * * * *", func() {
		notifyOneHourBeforeAppointment(client, messagingClient)
	})
	c.Start()

	return c
}
