package branch

import (
	"context"
	"math"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) ReviewAggregations(ctx context.Context, obj *prisma.Branch) (*gqlgen.ReviewAggregations, error) {
	approved := prisma.ReviewStatusApproved
	reviews, err := r.Prisma.Reviews(&prisma.ReviewsParams{
		Where: &prisma.ReviewWhereInput{
			Appointment: &prisma.AppointmentWhereInput{
				Branch: &prisma.BranchWhereInput{
					ID: &obj.ID,
				},
			},
			Status: &approved,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	sumRating := 0.0
	total := len(reviews)
	countPerStar := make([]int, 5)

	for _, review := range reviews {
		intStars := int(math.Ceil(review.Stars))
		countPerStar[intStars-1]++
		sumRating += review.Stars
	}

	if total > 0 {
		sumRating /= float64(total)
	}

	return &gqlgen.ReviewAggregations{
		TotalCount:    total,
		AverageRating: sumRating,
		CountPerStar:  countPerStar,
	}, nil
}
