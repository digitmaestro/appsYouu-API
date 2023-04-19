package appointment

import (
	"context"

	"firebase.google.com/go/messaging"

	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type Appointment struct {
	Prisma          *prisma.Client
	MessagingClient *messaging.Client
}

func (r *Appointment) MayUserModifyAppointment(ctx context.Context, appointmentID string) (bool, error) {
	user, err := sessctx.User(ctx)

	if err != nil {
		return false, err
	}

	switch user.Type {
	case prisma.UserTypeAdministrator:
		return true, nil
	case prisma.UserTypeCustomer:
		appointments, err := r.Prisma.Appointments(&prisma.AppointmentsParams{
			Where: &prisma.AppointmentWhereInput{
				ID: &appointmentID,
				Customer: &prisma.UserWhereInput{
					ID: &user.ID,
				},
			},
		}).Exec(ctx)

		if err != nil {
			return false, err
		}

		if len(appointments) == 0 {
			return false, gqlerrors.NewPermissionError("customer can only update his own appointments")
		}

	case prisma.UserTypeEmployee:
		appointments, err := r.Prisma.Appointments(&prisma.AppointmentsParams{
			Where: &prisma.AppointmentWhereInput{
				ID: &appointmentID,
				Employee: &prisma.UserWhereInput{
					ID: &user.ID,
				},
			},
		}).Exec(ctx)

		if err != nil {
			return false, err
		}

		if len(appointments) == 0 {
			return false, gqlerrors.NewPermissionError("employee can only update his own appointments")
		}

	case prisma.UserTypeManager:
		appointments, err := r.Prisma.Appointments(&prisma.AppointmentsParams{
			Where: &prisma.AppointmentWhereInput{
				ID: &appointmentID,
				Branch: &prisma.BranchWhereInput{
					Company: &prisma.CompanyWhereInput{
						UsersSome: &prisma.UserWhereInput{
							ID: &user.ID,
						},
					},
				},
			},
		}).Exec(ctx)

		if err != nil {
			return false, err
		}

		if len(appointments) == 0 {
			return false, gqlerrors.NewPermissionError("manager can only update appointments of his company branches")
		}
	default:
		return false, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	return true, nil
}

func New(client *prisma.Client, messagingClient *messaging.Client) *Appointment {
	return &Appointment{
		Prisma:          client,
		MessagingClient: messagingClient,
	}
}
