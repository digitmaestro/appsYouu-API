package product

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Product) ViewerReview(ctx context.Context, obj *prisma.Product) (*prisma.ProductReview, error) {
	viewer, err := sessctx.User(ctx)

	if obj.Deleted || err != nil {
		return nil, nil
	}

	reviewType := prisma.ReviewTypeProduct

	reviews, err := r.Prisma.Reviews(&prisma.ReviewsParams{
		Where: &prisma.ReviewWhereInput{
			Product: &prisma.ProductWhereInput{
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

	return reviews[0].Convert().(*prisma.ProductReview), err
}
