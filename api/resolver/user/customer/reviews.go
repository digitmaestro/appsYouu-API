package customer

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Customer) Reviews(ctx context.Context, obj *prisma.Customer) (*gqlgen.CustomerReviewConnection, error) {
	viewer, err := sessctx.User(ctx)

	nodes := []prisma.IReview{}

	if obj.Deleted || err != nil {
		return &gqlgen.CustomerReviewConnection{
			Nodes: nodes,
		}, nil
	}

	viewerCompany := sessctx.Company(ctx)

	where := &prisma.ReviewWhereInput{
		Customer: &prisma.UserWhereInput{
			ID: &obj.ID,
		},
	}

	if viewer.Type == prisma.UserTypeManager {
		where.Customer.Company = &prisma.CompanyWhereInput{
			ID: &viewerCompany,
		}
	} else if viewer.Type != prisma.UserTypeAdministrator {
		allowSharing := true
		approved := prisma.ReviewStatusApproved
		where.Customer.AllowReviewSharing = &allowSharing
		where.Status = &approved
	}

	reviews, err := r.Prisma.Reviews(&prisma.ReviewsParams{
		Where: where,
	}).Exec(ctx)

	if err != nil {
		return &gqlgen.CustomerReviewConnection{
			Nodes: nodes,
		}, err
	}

	for _, review := range reviews {
		clone := review
		nodes = append(nodes, clone.Convert())
	}

	return &gqlgen.CustomerReviewConnection{
		Nodes: nodes,
	}, err
}
