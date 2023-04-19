package manager

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/appointment"
	"github.com/steebchen/keskin-api/gqlgen"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Manager) Appointments(
	ctx context.Context,
	manager *prisma.Manager,
	input *gqlgen.AppointmentInput,
) (*gqlgen.AppointmentConnection, error) {
	nodes := []*prisma.Appointment{}

	if !manager.Deleted {
		if input == nil {
			input = &gqlgen.AppointmentInput{}
		}

		appointments, err := r.Prisma.Appointments(&prisma.AppointmentsParams{
			Where: &prisma.AppointmentWhereInput{
				Employee: &prisma.UserWhereInput{
					ID: &manager.ID,
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
