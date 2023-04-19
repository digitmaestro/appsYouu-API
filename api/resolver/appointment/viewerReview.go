package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) ViewerReview(ctx context.Context, obj *prisma.Appointment) (*prisma.AppointmentReview, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, nil
	}

	reviewType := prisma.ReviewTypeAppointment

	reviews, err := r.Prisma.Reviews(&prisma.ReviewsParams{
		Where: &prisma.ReviewWhereInput{
			Appointment: &prisma.AppointmentWhereInput{
				ID: &obj.ID,
			},
			Customer: &prisma.UserWhereInput{
				ID: &viewer.ID,
			},
			Type: &reviewType,
		},
	}).Exec(ctx)

	if err == prisma.ErrNoResult {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if len(reviews) == 0 {
		return nil, nil
	}

	return reviews[0].Convert().(*prisma.AppointmentReview), err
}
