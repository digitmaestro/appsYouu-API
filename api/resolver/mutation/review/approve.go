package review

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ReviewMutation) ApproveReview(
	ctx context.Context,
	id string,
	status prisma.ReviewStatus,
) (*gqlgen.ApproveReviewPayload, error) {
	company, err := r.Prisma.Review(prisma.ReviewWhereUniqueInput{
		ID: &id,
	}).Customer().Company().Exec(ctx)

	if err != nil {
		return nil, err
	}

	customer, err := r.Prisma.Review(prisma.ReviewWhereUniqueInput{
		ID: &id,
	}).Customer().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if customer.Deleted {
		return nil, gqlerrors.NewPermissionError("Customer is deleted")
	}

	if err := permissions.CanAccessCompany(ctx, company.ID, r.Prisma, []prisma.UserType{prisma.UserTypeManager}); err != nil {
		return nil, err
	}

	review, err := r.Prisma.UpdateReview(prisma.ReviewUpdateParams{
		Where: prisma.ReviewWhereUniqueInput{
			ID: &id,
		},
		Data: prisma.ReviewUpdateInput{
			Status: &status,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.ApproveReviewPayload{
		Review: review.Convert(),
	}, nil
}
