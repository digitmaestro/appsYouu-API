package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) Products(ctx context.Context, obj *prisma.Appointment) ([]*gqlgen.AppointmentProduct, error) {
	result := []*gqlgen.AppointmentProduct{}

	productLinks, err := r.Prisma.AppointmentProductLinks(&prisma.AppointmentProductLinksParams{
		Where: &prisma.AppointmentProductLinkWhereInput{
			Appointment: &prisma.AppointmentWhereInput{
				ID: &obj.ID,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	for _, productLink := range productLinks {
		product, err := r.Prisma.AppointmentProductLink(prisma.AppointmentProductLinkWhereUniqueInput{
			ID: &productLink.ID,
		}).Product().Exec(ctx)

		if err != nil {
			return nil, err
		}

		clone := product
		result = append(result, &gqlgen.AppointmentProduct{
			Item:  clone,
			Count: int(productLink.Count),
		})
	}

	return result, nil
}
