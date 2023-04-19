package appointment

import (
	"context"
	"fmt"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/api/resolver/mutation/email_template"
	"github.com/steebchen/keskin-api/api/resolver/mutation/notification"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

var statusApproved = prisma.AppointmentStatusApproved

func (r *Appointment) StaffApproveAppointment(
	ctx context.Context,
	input gqlgen.StaffApproveAppointmentInput,
	language *string,
) (*gqlgen.StaffApproveAppointmentPayload, error) {
	branch, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &input.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	t := []prisma.UserType{prisma.UserTypeEmployee, prisma.UserTypeManager}
	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, t); err != nil {
		return nil, err
	}

	if input.Patch == nil {
		input.Patch = &gqlgen.StaffApproveAppointmentPatch{}
	}

	appointment, err := r.Prisma.UpdateAppointment(prisma.AppointmentUpdateParams{
		Where: prisma.AppointmentWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.AppointmentUpdateInput{
			Status: &statusApproved,
			Desc:   i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Desc),
			Note:   input.Patch.Note,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	customer, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &input.ID,
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

	return &gqlgen.StaffApproveAppointmentPayload{
		Appointment: appointment,
	}, nil
}
