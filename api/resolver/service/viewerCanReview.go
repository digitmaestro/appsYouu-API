package service

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Service) ViewerCanReview(ctx context.Context, obj *prisma.Service) (*bool, error) {
	viewer, err := sessctx.User(ctx)

	result := false

	if obj.Deleted || err != nil {
		return &result, nil
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

	if err != nil && err != prisma.ErrNoResult {
		return &result, err
	}

	if len(reviews) > 0 {
		return &result, nil
	}

	services, err := r.Prisma.Services(&prisma.ServicesParams{
		Where: &prisma.ServiceWhereInput{
			ID: &obj.ID,
			Branch: &prisma.BranchWhereInput{
				Company: &prisma.CompanyWhereInput{
					UsersSome: &prisma.UserWhereInput{
						ID: &viewer.ID,
					},
				},
			},
		},
	}).Exec(ctx)

	if err != nil && err != prisma.ErrNoResult {
		return &result, err
	}

	if len(services) == 0 {
		return &result, nil
	}

	result = true

	return &result, nil
}
