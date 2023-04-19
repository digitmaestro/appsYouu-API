package product

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type ProductQuery struct {
	Prisma *prisma.Client
}

func (r *ProductQuery) Product(ctx context.Context, id string, language *string) (*prisma.Product, error) {
	return r.Prisma.Product(prisma.ProductWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)
}

func (r *ProductQuery) Products(ctx context.Context, input gqlgen.ProductInput, language *string) (*gqlgen.ProductConnection, error) {
	deleted := false
	where := &prisma.ProductWhereInput{
		Deleted: &deleted,
	}

	if input.Branch == nil {
		companyId := sessctx.CompanyWithFallback(ctx, r.Prisma, input.Company)

		where.Branch = &prisma.BranchWhereInput{
			Company: &prisma.CompanyWhereInput{
				ID: &companyId,
			},
		}
	} else {
		where.Branch = &prisma.BranchWhereInput{
			ID: input.Branch,
		}
	}

	products, err := r.Prisma.Products(&prisma.ProductsParams{
		Where:   where,
		OrderBy: assembleOrder(input.Order),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	nodes := []*prisma.Product{}

	for _, product := range products {
		clone := product
		nodes = append(nodes, &clone)
	}

	return &gqlgen.ProductConnection{
		Nodes: nodes,
	}, nil
}

func New(client *prisma.Client) *ProductQuery {
	return &ProductQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
