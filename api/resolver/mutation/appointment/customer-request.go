package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) CustomerRequestAppointment(
	ctx context.Context,
	input gqlgen.CustomerRequestAppointmentInput,
	language *string,
) (*gqlgen.CustomerRequestAppointmentPayload, error) {
	user, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	_, err = r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &input.Branch,
	}).Exec(ctx)

	if err != nil {
		return nil, gqlerrors.NewFormatNodeError(err, input.Branch)
	}

	if input.Employee != nil {
		staff, err := r.Prisma.User(prisma.UserWhereUniqueInput{
			ID: input.Employee,
		}).Exec(ctx)

		if err != nil {
			return nil, gqlerrors.NewFormatNodeError(err, *input.Employee)
		}

		if staff.Type != prisma.UserTypeManager && staff.Type != prisma.UserTypeEmployee {
			return nil, gqlerrors.NewValidationError("staff type is invalid", "InvalidUserType")
		}

		if staff.Deleted {
			return nil, gqlerrors.NewPermissionError("Staff is deleted")
		}
	}

	appointment, err := CreateAppointment(CreateAppointmentInput{
		Client:          r.Prisma,
		Context:         ctx,
		EmployeeID:      input.Employee,
		CustomerID:      user.ID,
		Branch:          input.Branch,
		Desc:            input.Data.Desc,
		Start:           input.Data.Start,
		ProductRequests: input.Data.Products,
		ServiceRequests: input.Data.Services,
		DefaultStatus:   prisma.AppointmentStatusRequested,
	})

	if err != nil {
		return nil, err
	}

	return &gqlgen.CustomerRequestAppointmentPayload{
		Appointment: appointment,
	}, nil
}
