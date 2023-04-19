package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) Employee(ctx context.Context, obj *prisma.Appointment) (*prisma.Employee, error) {
	employee, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &obj.ID,
	}).Employee().Exec(ctx)

	return &prisma.Employee{
		User: employee,
	}, err
}
