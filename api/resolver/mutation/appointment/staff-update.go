package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) StaffUpdateAppointment(
	ctx context.Context,
	input gqlgen.StaffUpdateAppointmentInput,
	language *string,
) (*gqlgen.StaffUpdateAppointmentPayload, error) {
	branch, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &input.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, gqlerrors.NewFormatNodeError(err, input.ID)
	}

	t := []prisma.UserType{prisma.UserTypeEmployee, prisma.UserTypeManager}
	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, t); err != nil {
		return nil, err
	}

	appointment, err := r.Prisma.UpdateAppointment(prisma.AppointmentUpdateParams{
		Where: prisma.AppointmentWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.AppointmentUpdateInput{
			Note: input.Patch.Note,
			Desc: i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Desc),
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.StaffUpdateAppointmentPayload{
		Appointment: appointment,
	}, nil
}
