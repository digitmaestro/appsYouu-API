package employee

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/appointment"
	"github.com/steebchen/keskin-api/gqlgen"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Employee) Appointments(
	ctx context.Context,
	employee *prisma.Employee,
	input *gqlgen.AppointmentInput,
) (*gqlgen.AppointmentConnection, error) {
	nodes := []*prisma.Appointment{}

	if !employee.Deleted {
		if input == nil {
			input = &gqlgen.AppointmentInput{}
		}

		appointments, err := r.Prisma.Appointments(&prisma.AppointmentsParams{
			Where: &prisma.AppointmentWhereInput{
				Employee: &prisma.UserWhereInput{
					ID: &employee.ID,
				},
				StatusIn: input.Status,
				StartGt:  gqlgen.TimeFilter(input.Start).Gt,
				StartGte: gqlgen.TimeFilter(input.Start).Gte,
				StartLt:  gqlgen.TimeFilter(input.Start).Lt,
				StartLte: gqlgen.TimeFilter(input.Start).Lte,
				EndGt:    gqlgen.TimeFilter(input.End).Gt,
				EndGte:   gqlgen.TimeFilter(input.End).Gte,
				EndLt:    gqlgen.TimeFilter(input.End).Lt,
				EndLte:   gqlgen.TimeFilter(input.End).Lte,
			},
			OrderBy: appointment.AssembleOrder(input.Order),
		}).Exec(ctx)

		if err != nil {
			return &gqlgen.AppointmentConnection{
				Nodes: nodes,
			}, err
		}

		for _, item := range appointments {
			clone := item
			nodes = append(nodes, &clone)
		}
	}

	return &gqlgen.AppointmentConnection{
		Nodes: nodes,
	}, nil
}
