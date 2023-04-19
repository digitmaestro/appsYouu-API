package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) ViewerCanReview(ctx context.Context, obj *prisma.Appointment) (*bool, error) {
	viewer, err := sessctx.User(ctx)

	result := false

	if err != nil {
		return &result, nil
	}

	review, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &obj.ID,
	}).Review().Exec(ctx)

	if err != nil && err != prisma.ErrNoResult {
		return &result, err
	}

	if review != nil {
		return &result, nil
	}

	customer, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &obj.ID,
	}).Customer().Exec(ctx)

	if err != nil {
		return &result, err
	}

	result = customer.ID == viewer.ID && obj.Status == prisma.AppointmentStatusApproved

	return &result, nil
}
