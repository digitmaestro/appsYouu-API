package appointment_review

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *AppointmentReview) Customer(ctx context.Context, obj *prisma.AppointmentReview) (*prisma.Customer, error) {
	customer, err := r.Prisma.Review(prisma.ReviewWhereUniqueInput{
		ID: &obj.ID,
	}).Customer().Exec(ctx)

	return &prisma.Customer{
		User: customer,
	}, err
}
