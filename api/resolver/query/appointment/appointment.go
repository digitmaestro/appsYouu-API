package appointment

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type AppointmentQuery struct {
	Prisma *prisma.Client
}

func (r *AppointmentQuery) Appointment(ctx context.Context, id string, language *string) (*prisma.Appointment, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	var appointments []prisma.Appointment

	switch viewer.Type {
	case prisma.UserTypeAdministrator:
		appointments, err = r.Prisma.Appointments(&prisma.AppointmentsParams{
			Where: &prisma.AppointmentWhereInput{
				ID: &id,
			},
		}).Exec(ctx)
	case prisma.UserTypeManager:
		appointments, err = r.Prisma.Appointments(&prisma.AppointmentsParams{
			Where: &prisma.AppointmentWhereInput{
				ID: &id,
				Branch: &prisma.BranchWhereInput{
					Company: &prisma.CompanyWhereInput{
						UsersSome: &prisma.UserWhereInput{
							ID: &viewer.ID,
						},
					},
				},
			},
		}).Exec(ctx)
	case prisma.UserTypeCustomer:
		appointments, err = r.Prisma.Appointments(&prisma.AppointmentsParams{
			Where: &prisma.AppointmentWhereInput{
				ID: &id,
				Customer: &prisma.UserWhereInput{
					ID: &viewer.ID,
				},
			},
		}).Exec(ctx)
	case prisma.UserTypeEmployee:
		appointments, err = r.Prisma.Appointments(&prisma.AppointmentsParams{
			Where: &prisma.AppointmentWhereInput{
				ID: &id,
				Employee: &prisma.UserWhereInput{
					ID: &viewer.ID,
				},
			},
		}).Exec(ctx)
	default:
		return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	if err != nil {
		return nil, err
	}

	if len(appointments) > 0 {
		return &appointments[0], nil
	}

	return nil, gqlerrors.NewNotFoundError(id)
}

func New(client *prisma.Client) *AppointmentQuery {
	return &AppointmentQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
