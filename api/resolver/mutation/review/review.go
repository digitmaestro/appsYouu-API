package review

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

var allowedTypes = []prisma.UserType{
	prisma.UserTypeManager,
}

type ReviewMutation struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *ReviewMutation {
	return &ReviewMutation{
		Prisma: client,
	}
}

func CreateReviewInput(patch *gqlgen.CreateReviewData) prisma.ReviewCreateInput {
	return prisma.ReviewCreateInput{
		Stars: patch.Stars,
		Title: patch.Title,
		Text:  patch.Text,
	}
}

func CreateProductReview(
	ctx context.Context,
	prismaClient *prisma.Client,
	input gqlgen.CreateProductReviewInput,
) (*prisma.Review, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	reviews, err := prismaClient.Reviews(&prisma.ReviewsParams{
		Where: &prisma.ReviewWhereInput{
			Customer: &prisma.UserWhereInput{
				ID: &viewer.ID,
			},
			Product: &prisma.ProductWhereInput{
				ID: &input.Product,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(reviews) > 0 {
		return nil, gqlerrors.NewPermissionError("only one review per customer allowed")
	}

	products, err := prismaClient.Products(&prisma.ProductsParams{
		Where: &prisma.ProductWhereInput{
			ID: &input.Product,
			Branch: &prisma.BranchWhereInput{
				Company: &prisma.CompanyWhereInput{
					UsersSome: &prisma.UserWhereInput{
						ID: &viewer.ID,
					},
				},
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, gqlerrors.NewPermissionError("user is not allowed to review this product")
	}

	return prismaClient.CreateReview(prisma.ReviewCreateInput{
		Stars: input.Review.Stars,
		Title: input.Review.Title,
		Text:  input.Review.Text,
		Product: &prisma.ProductCreateOneWithoutReviewsInput{
			Connect: &prisma.ProductWhereUniqueInput{
				ID: &input.Product,
			},
		},
		Type: prisma.ReviewTypeProduct,
		Customer: prisma.UserCreateOneWithoutReviewsInput{
			Connect: &prisma.UserWhereUniqueInput{
				ID: &viewer.ID,
			},
		},
	}).Exec(ctx)
}

func CreateAppointmentReview(
	ctx context.Context,
	prismaClient *prisma.Client,
	input gqlgen.CreateAppointmentReviewInput,
) (*prisma.Review, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	reviews, err := prismaClient.Reviews(&prisma.ReviewsParams{
		Where: &prisma.ReviewWhereInput{
			Customer: &prisma.UserWhereInput{
				ID: &viewer.ID,
			},
			Appointment: &prisma.AppointmentWhereInput{
				ID: &input.Appointment,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(reviews) > 0 {
		return nil, gqlerrors.NewPermissionError("only one review per customer allowed")
	}

	appovedStatus := prisma.AppointmentStatusApproved

	appointments, err := prismaClient.Appointments(&prisma.AppointmentsParams{
		Where: &prisma.AppointmentWhereInput{
			Status: &appovedStatus,
			Customer: &prisma.UserWhereInput{
				ID: &viewer.ID,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(appointments) == 0 {
		return nil, gqlerrors.NewPermissionError("user is not allowed to review this appointment")
	}

	return prismaClient.CreateReview(prisma.ReviewCreateInput{
		Stars: input.Review.Stars,
		Title: input.Review.Title,
		Text:  input.Review.Text,
		Appointment: &prisma.AppointmentCreateOneWithoutReviewInput{
			Connect: &prisma.AppointmentWhereUniqueInput{
				ID: &input.Appointment,
			},
		},
		Type: prisma.ReviewTypeAppointment,
		Customer: prisma.UserCreateOneWithoutReviewsInput{
			Connect: &prisma.UserWhereUniqueInput{
				ID: &viewer.ID,
			},
		},
	}).Exec(ctx)
}

func CreateServiceReview(
	ctx context.Context,
	prismaClient *prisma.Client,
	input gqlgen.CreateServiceReviewInput,
) (*prisma.Review, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	reviews, err := prismaClient.Reviews(&prisma.ReviewsParams{
		Where: &prisma.ReviewWhereInput{
			Customer: &prisma.UserWhereInput{
				ID: &viewer.ID,
			},
			Service: &prisma.ServiceWhereInput{
				ID: &input.Service,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(reviews) > 0 {
		return nil, gqlerrors.NewPermissionError("only one review per customer allowed")
	}

	services, err := prismaClient.Services(&prisma.ServicesParams{
		Where: &prisma.ServiceWhereInput{
			ID: &input.Service,
			Branch: &prisma.BranchWhereInput{
				Company: &prisma.CompanyWhereInput{
					UsersSome: &prisma.UserWhereInput{
						ID: &viewer.ID,
					},
				},
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, gqlerrors.NewPermissionError("user is not allowed to review this service")
	}

	return prismaClient.CreateReview(prisma.ReviewCreateInput{
		Stars: input.Review.Stars,
		Title: input.Review.Title,
		Text:  input.Review.Text,
		Service: &prisma.ServiceCreateOneWithoutReviewsInput{
			Connect: &prisma.ServiceWhereUniqueInput{
				ID: &input.Service,
			},
		},
		Type: prisma.ReviewTypeService,
		Customer: prisma.UserCreateOneWithoutReviewsInput{
			Connect: &prisma.UserWhereUniqueInput{
				ID: &viewer.ID,
			},
		},
	}).Exec(ctx)
}
