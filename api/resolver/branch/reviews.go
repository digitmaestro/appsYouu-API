package branch

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) Reviews(ctx context.Context, obj *prisma.Branch) (*gqlgen.AppointmentReviewConnection, error) {
	viewer, err := sessctx.User(ctx)

	nodes := []*prisma.AppointmentReview{}

	if err != nil {
		return &gqlgen.AppointmentReviewConnection{
			Nodes: nodes,
		}, nil
	}

	viewerCompany := sessctx.Company(ctx)

	where := &prisma.ReviewWhereInput{
		Appointment: &prisma.AppointmentWhereInput{
			Branch: &prisma.BranchWhereInput{
				ID: &obj.ID,
			},
		},
	}

	if viewer.Type != prisma.UserTypeAdministrator {
		where.Or = []prisma.ReviewWhereInput{{
			Customer: &prisma.UserWhereInput{
				ID: &viewer.ID,
			},
		}}

		if viewer.Type == prisma.UserTypeManager {
			where.Or = append(where.Or, prisma.ReviewWhereInput{
				Customer: &prisma.UserWhereInput{
					Company: &prisma.CompanyWhereInput{
						ID: &viewerCompany,
					},
				},
			})
		} else if viewer.Type != prisma.UserTypeAdministrator {
			allowSharing := true
			approved := prisma.ReviewStatusApproved
			where.Or = append(where.Or, prisma.ReviewWhereInput{
				Customer: &prisma.UserWhereInput{
					AllowReviewSharing: &allowSharing,
				},
				Status: &approved,
			})
		}
	}

	reviews, err := r.Prisma.Reviews(&prisma.ReviewsParams{
		Where: where,
	}).Exec(ctx)

	if err != nil {
		return &gqlgen.AppointmentReviewConnection{
			Nodes: nodes,
		}, err
	}

	for _, review := range reviews {
		clone := review
		nodes = append(nodes, clone.Convert().(*prisma.AppointmentReview))
	}

	return &gqlgen.AppointmentReviewConnection{
		Nodes: nodes,
	}, err
}
