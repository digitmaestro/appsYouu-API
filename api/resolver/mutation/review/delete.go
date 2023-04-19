package review

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ReviewMutation) DeleteReview(
	ctx context.Context,
	id string,
) (*gqlgen.DeleteReviewPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	customer, err := r.Prisma.Review(prisma.ReviewWhereUniqueInput{
		ID: &id,
	}).Customer().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if viewer.Deleted || (viewer.ID != customer.ID && viewer.Type != prisma.UserTypeManager && viewer.Type != prisma.UserTypeAdministrator) {
		return nil, gqlerrors.NewPermissionError("user is not allowed to access this review")
	}

	review, err := r.Prisma.Review(prisma.ReviewWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	var branch *prisma.Branch = nil

	if review.Type == prisma.ReviewTypeAppointment {
		branch, err = r.Prisma.Review(prisma.ReviewWhereUniqueInput{
			ID: &id,
		}).Appointment().Branch().Exec(ctx)
	} else if review.Type == prisma.ReviewTypeProduct {
		branch, err = r.Prisma.Review(prisma.ReviewWhereUniqueInput{
			ID: &id,
		}).Product().Branch().Exec(ctx)
	} else if review.Type == prisma.ReviewTypeService {
		branch, err = r.Prisma.Review(prisma.ReviewWhereUniqueInput{
			ID: &id,
		}).Service().Branch().Exec(ctx)
	}

	if err != nil {
		return nil, err
	}

	if branch == nil {
		return nil, gqlerrors.NewPermissionError("user is not allowed to access this branch")
	}

	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, allowedTypes); err != nil && viewer.ID != customer.ID {
		return nil, err
	}

	review, err = r.Prisma.DeleteReview(prisma.ReviewWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.DeleteReviewPayload{
		Review: review.Convert(),
	}, nil
}
