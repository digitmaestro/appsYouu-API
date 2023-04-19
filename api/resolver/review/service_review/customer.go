package service_review

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *ServiceReview) Customer(ctx context.Context, obj *prisma.ServiceReview) (*prisma.Customer, error) {
	customer, err := r.Prisma.Review(prisma.ReviewWhereUniqueInput{
		ID: &obj.ID,
	}).Customer().Exec(ctx)

	return &prisma.Customer{
		User: customer,
	}, err
}
