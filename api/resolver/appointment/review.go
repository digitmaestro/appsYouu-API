package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) Review(ctx context.Context, obj *prisma.Appointment) (*prisma.AppointmentReview, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, nil
	}

	viewerCompany := sessctx.Company(ctx)

	company, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &obj.ID,
	}).Branch().Company().Exec(ctx)

	if err != nil {
		return nil, nil
	}

	customer, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &obj.ID,
	}).Customer().Exec(ctx)

	if err != nil {
		return nil, nil
	}

	review, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &obj.ID,
	}).Review().Exec(ctx)

	if err == prisma.ErrNoResult {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var appointmentReview *prisma.AppointmentReview = nil

	if (review.Status == prisma.ReviewStatusApproved && customer.AllowReviewSharing) || viewer.Type == prisma.UserTypeAdministrator || (viewer.Type == prisma.UserTypeManager && viewerCompany == company.ID) || customer.ID == viewer.ID {
		appointmentReview = review.Convert().(*prisma.AppointmentReview)
	}

	return appointmentReview, nil
}
