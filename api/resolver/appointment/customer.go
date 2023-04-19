package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) Customer(ctx context.Context, obj *prisma.Appointment) (*prisma.Customer, error) {
	customer, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &obj.ID,
	}).Customer().Exec(ctx)

	return &prisma.Customer{
		User: customer,
	}, err
}
