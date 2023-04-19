package service_review

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *ServiceReview) Service(ctx context.Context, obj *prisma.ServiceReview) (*prisma.Service, error) {
	service, err := r.Prisma.Review(prisma.ReviewWhereUniqueInput{
		ID: &obj.ID,
	}).Service().Exec(ctx)

	return service, err
}
