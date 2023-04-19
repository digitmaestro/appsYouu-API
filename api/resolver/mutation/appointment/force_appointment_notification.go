package appointment

import (
	"context"
	"fmt"

	"github.com/steebchen/keskin-api/api/resolver/mutation/email_template"
	"github.com/steebchen/keskin-api/api/resolver/mutation/notification"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) ForceAppointmentNotification(
	ctx context.Context,
	input *gqlgen.ForceAppointmentNotificationInput,
) (*gqlgen.NotificationPayload, error) {
	authorized, err := r.MayUserModifyAppointment(ctx, input.Appointment)

	if err != nil {
		return nil, err
	}

	if !authorized {
		return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	appointment, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &input.Appointment,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	customer, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &input.Appointment,
	}).Customer().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if customer.Deleted {
		return nil, gqlerrors.NewPermissionError("Customer is deleted")
	}

	branch, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &input.Appointment,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	appointmentDate := i18n.FormatDate(ctx, appointment.Start)
	appointmentTime := i18n.FormatTime(ctx, appointment.Start)

	notificationContext := context.Background()
	notificationContext = sessctx.SetLanguage(notificationContext, customer.Language)

	response, err := notification.Send(
		r.Prisma,
		r.MessagingClient,
		notificationContext,
		customer.ID,
		i18n.Language(notificationContext)["APPOINTMENT_REMINDER_TITLE"],
		fmt.Sprintf(
			i18n.Language(notificationContext)["APPOINTMENT_REMINDER_TEXT"],
			appointmentDate,
			appointmentTime,
		),
	)

	go email_template.SendEmailTemplate(
		notificationContext,
		r.Prisma,
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

	return &gqlgen.NotificationPayload{
		UserID: customer.ID,
		Result: response,
	}, nil
}
