package review

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type ReviewQuery struct {
	Prisma *prisma.Client
}

func (r *ReviewQuery) Reviews(ctx context.Context, input gqlgen.ReviewInput, language *string) (*gqlgen.ReviewConnection, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	viewerCompany := sessctx.CompanyWithFallback(ctx, r.Prisma, input.Company)

	nodes := []prisma.IReview{}

	where := &prisma.ReviewWhereInput{
		Customer: &prisma.UserWhereInput{},
	}

	if viewer.Type == prisma.UserTypeManager ||
		viewer.Type == prisma.UserTypeAdministrator {
		if input.Customer != nil {
			where.Customer.ID = input.Customer
		} else {
			where.Customer.Company = &prisma.CompanyWhereInput{
				ID: &viewerCompany,
			}
		}
	} else if viewer.Type == prisma.UserTypeCustomer {
		where.Customer.ID = &viewer.ID
	} else {
		return &gqlgen.ReviewConnection{
			Nodes: nodes,
		}, nil
	}

	if input.Status != nil && len(input.Status) > 0 {
		where.StatusIn = input.Status
	}

	if input.Type != nil && len(input.Type) > 0 {
		where.TypeIn = input.Type
	}

	reviews, err := r.Prisma.Reviews(&prisma.ReviewsParams{
		Where:   where,
		OrderBy: AssembleOrder(input.Order),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	for _, review := range reviews {
		clone := review
		nodes = append(nodes, clone.Convert())
	}

	return &gqlgen.ReviewConnection{
		Nodes: nodes,
	}, err
}

func New(client *prisma.Client) *ReviewQuery {
	return &ReviewQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
