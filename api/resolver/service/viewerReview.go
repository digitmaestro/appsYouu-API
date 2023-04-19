package service

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Service) ViewerReview(ctx context.Context, obj *prisma.Service) (*prisma.ServiceReview, error) {
	viewer, err := sessctx.User(ctx)

	if obj.Deleted || err != nil {
		return nil, nil
	}

	reviewType := prisma.ReviewTypeService

	reviews, err := r.Prisma.Reviews(&prisma.ReviewsParams{
		Where: &prisma.ReviewWhereInput{
			Service: &prisma.ServiceWhereInput{
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

	return reviews[0].Convert().(*prisma.ServiceReview), err
}
