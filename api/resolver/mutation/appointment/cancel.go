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

var statusCanceled = prisma.AppointmentStatusCanceled

func (r *Appointment) CancelAppointment(
	ctx context.Context,
	input gqlgen.CancelAppointmentInput,
	language *string,
) (*gqlgen.CancelAppointmentPayload, error) {
	authorized, err := r.MayUserModifyAppointment(ctx, input.ID)

	if err != nil {
		return nil, err
	}

	if !authorized {
		return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	customer, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &input.ID,
	}).Customer().Exec(ctx)

	if err != nil {
		return nil, err
	}

	appointment, err := r.Prisma.UpdateAppointment(prisma.AppointmentUpdateParams{
		Where: prisma.AppointmentWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.AppointmentUpdateInput{
			Status: &statusCanceled,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	branch, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &input.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	user, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	appointmentDate := i18n.FormatDate(ctx, appointment.Start)
	appointmentTime := i18n.FormatTime(ctx, appointment.Start)

	notificationContext := context.Background()
	notificationContext = sessctx.SetLanguage(notificationContext, customer.Language)

	if user.ID == customer.ID {
		go notification.Send(
			r.Prisma,
			r.MessagingClient,
			notificationContext,
			customer.ID,
			i18n.Language(notificationContext)["APPOINTMENT_CANCELED_BY_USER_TITLE"],
			i18n.Language(notificationContext)["APPOINTMENT_CANCELED_BY_USER_TEXT"],
		)

		go email_template.SendEmailTemplate(
			notificationContext,
			r.Prisma,
			"appointmentCanceledByUser",
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
	} else {
		go notification.Send(
			r.Prisma,
			r.MessagingClient,
			notificationContext,
			customer.ID,
			i18n.Language(notificationContext)["APPOINTMENT_CANCELED_TITLE"],
			fmt.Sprintf(
				i18n.Language(notificationContext)["APPOINTMENT_CANCELED_TEXT"],
				appointmentDate,
				appointmentTime,
			),
		)

		go email_template.SendEmailTemplate(
			notificationContext,
			r.Prisma,
			"appointmentCanceled",
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

	return &gqlgen.CancelAppointmentPayload{
		Appointment: appointment,
	}, nil
}
