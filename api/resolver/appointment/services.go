package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) Services(ctx context.Context, obj *prisma.Appointment) ([]*prisma.Service, error) {
	result := []*prisma.Service{}

	serviceLinks, err := r.Prisma.AppointmentServiceLinks(&prisma.AppointmentServiceLinksParams{
		Where: &prisma.AppointmentServiceLinkWhereInput{
			Appointment: &prisma.AppointmentWhereInput{
				ID: &obj.ID,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	for _, serviceLink := range serviceLinks {
		service, err := r.Prisma.AppointmentServiceLink(prisma.AppointmentServiceLinkWhereUniqueInput{
			ID: &serviceLink.ID,
		}).Service().Exec(ctx)

		if err != nil {
			return nil, err
		}

		clone := service
		result = append(result, clone)
	}

	return result, nil
}
