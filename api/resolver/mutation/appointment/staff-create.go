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

func (r *Appointment) StaffCreateAppointment(
	ctx context.Context,
	input gqlgen.StaffCreateAppointmentInput,
	language *string,
) (*gqlgen.StaffCreateAppointmentPayload, error) {
	employeeID := input.Employee

	if employeeID == nil {
		user, err := sessctx.User(ctx)

		if err != nil {
			return nil, err
		}

		employeeID = &user.ID
	}

	appointment, err := CreateAppointment(CreateAppointmentInput{
		Client:          r.Prisma,
		Context:         ctx,
		EmployeeID:      employeeID,
		CustomerID:      input.Customer,
		Branch:          input.Branch,
		Desc:            input.Data.Desc,
		Start:           input.Data.Start,
		ProductRequests: input.Data.Products,
		ServiceRequests: input.Data.Services,
		DefaultStatus:   prisma.AppointmentStatusApproved,
	})

	if err != nil {
		return nil, err
	}

	customer, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &appointment.ID,
	}).Customer().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if customer.Deleted {
		return nil, gqlerrors.NewPermissionError("Customer is deleted")
	}

	appointmentDate := i18n.FormatDate(ctx, appointment.Start)
	appointmentTime := i18n.FormatTime(ctx, appointment.Start)

	notificationContext := context.Background()
	notificationContext = sessctx.SetLanguage(notificationContext, customer.Language)

	go notification.Send(
		r.Prisma,
		r.MessagingClient,
		notificationContext,
		customer.ID,
		i18n.Language(notificationContext)["APPOINTMENT_APPROVED_TITLE"],
		fmt.Sprintf(
			i18n.Language(notificationContext)["APPOINTMENT_APPROVED_TEXT"],
			appointmentDate,
			appointmentTime,
		),
	)

	go email_template.SendEmailTemplate(
		notificationContext,
		r.Prisma,
		"appointmentConfirmed",
		input.Branch,
		customer.Email,
		customer.Gender,
		customer.LastName,
		customer.FirstName,
		&appointmentDate,
		&appointmentTime,
		nil,
		nil,
	)

	return &gqlgen.StaffCreateAppointmentPayload{
		Appointment: appointment,
	}, nil
}
