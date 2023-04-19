package product

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Product) ViewerCanReview(ctx context.Context, obj *prisma.Product) (*bool, error) {
	viewer, err := sessctx.User(ctx)

	result := false

	if obj.Deleted || err != nil {
		return &result, nil
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

	if err != nil && err != prisma.ErrNoResult {
		return &result, err
	}

	if len(reviews) > 0 {
		return &result, nil
	}

	products, err := r.Prisma.Products(&prisma.ProductsParams{
		Where: &prisma.ProductWhereInput{
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

	if len(products) == 0 {
		return &result, nil
	}

	result = true

	return &result, nil
}
