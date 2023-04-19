package appointment_review

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *AppointmentReview) Appointment(ctx context.Context, obj *prisma.AppointmentReview) (*prisma.Appointment, error) {
	appointment, err := r.Prisma.Review(prisma.ReviewWhereUniqueInput{
		ID: &obj.ID,
	}).Appointment().Exec(ctx)

	return appointment, err
}
