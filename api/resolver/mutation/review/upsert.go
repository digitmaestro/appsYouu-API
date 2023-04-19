package review

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ReviewMutation) UpsertReview(
	ctx context.Context,
	input gqlgen.UpsertReviewInput,
) (*gqlgen.UpsertReviewPayload, error) {
	if input.Product == nil && input.Service == nil && input.Appointment == nil && input.Review == nil {
		return nil, gqlerrors.NewValidationError("At least one of Product, Service, Appointment or Review is required", "MissingData")
	}

	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if viewer.Deleted {
		return nil, gqlerrors.NewPermissionError("user is not allowed to access this review")
	}

	var existingReview *prisma.Review = nil

	existingPreviewsWhere := &prisma.ReviewWhereInput{
		Customer: &prisma.UserWhereInput{
			ID: &viewer.ID,
		},
	}

	if input.Product != nil {
		existingPreviewsWhere.Product = &prisma.ProductWhereInput{
			ID: input.Product,
		}
	} else if input.Service != nil {
		existingPreviewsWhere.Service = &prisma.ServiceWhereInput{
			ID: input.Service,
		}
	} else if input.Appointment != nil {
		existingPreviewsWhere.Appointment = &prisma.AppointmentWhereInput{
			ID: input.Appointment,
		}
	} else if input.Review != nil {
		existingPreviewsWhere.ID = input.Review
	}

	existingReviews, err := r.Prisma.Reviews(&prisma.ReviewsParams{
		Where: existingPreviewsWhere,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(existingReviews) > 0 {
		existingReview = &existingReviews[0]
	}

	var review *prisma.Review = nil

	if existingReview == nil {
		if input.Review != nil {
			return nil, gqlerrors.NewPermissionError("user is not allowed to access this review")
		}

		if input.Data.Stars == nil {
			return nil, gqlerrors.NewValidationError("Stars is required", "MissingData")
		}

		if input.Data.Title == nil {
			return nil, gqlerrors.NewValidationError("Title is required", "MissingData")
		}

		if input.Data.Text == nil {
			return nil, gqlerrors.NewValidationError("Text is required", "MissingData")
		}

		create := gqlgen.CreateReviewData{
			Stars: *input.Data.Stars,
			Title: *input.Data.Title,
			Text:  *input.Data.Text,
		}

		if input.Product != nil {
			review, err = CreateProductReview(
				ctx,
				r.Prisma,
				gqlgen.CreateProductReviewInput{
					Product: *input.Product,
					Review:  &create,
				},
			)
		} else if input.Service != nil {
			review, err = CreateServiceReview(
				ctx,
				r.Prisma,
				gqlgen.CreateServiceReviewInput{
					Service: *input.Service,
					Review:  &create,
				},
			)
		} else if input.Appointment != nil {
			review, err = CreateAppointmentReview(
				ctx,
				r.Prisma,
				gqlgen.CreateAppointmentReviewInput{
					Appointment: *input.Appointment,
					Review:      &create,
				},
			)
		}
	} else {
		pending := prisma.ReviewStatusPending

		review, err = r.Prisma.UpdateReview(prisma.ReviewUpdateParams{
			Where: prisma.ReviewWhereUniqueInput{
				ID: &existingReview.ID,
			},
			Data: prisma.ReviewUpdateInput{
				Stars:  input.Data.Stars,
				Title:  input.Data.Title,
				Text:   input.Data.Text,
				Status: &pending,
			},
		}).Exec(ctx)
	}

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpsertReviewPayload{
		Review: review.Convert(),
	}, nil
}
