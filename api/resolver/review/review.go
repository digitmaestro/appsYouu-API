package review

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/api/resolver/review/appointment_review"
	"github.com/steebchen/keskin-api/api/resolver/review/product_review"
	"github.com/steebchen/keskin-api/api/resolver/review/service_review"
	"github.com/steebchen/keskin-api/prisma"
)

type Review struct {
	Prisma                    *prisma.Client
	ProductReviewResolver     *product_review.ProductReview
	ServiceReviewResolver     *service_review.ServiceReview
	AppointmentReviewResolver *appointment_review.AppointmentReview
}

func New(
	client *prisma.Client,
	productReview *product_review.ProductReview,
	serviceReview *service_review.ServiceReview,
	appointmentReview *appointment_review.AppointmentReview,
) *Review {
	return &Review{
		Prisma:                    client,
		ProductReviewResolver:     productReview,
		ServiceReviewResolver:     serviceReview,
		AppointmentReviewResolver: appointmentReview,
	}
}

var ProviderSet = wire.NewSet(
	New,
	product_review.New,
	service_review.New,
	appointment_review.New,
)
