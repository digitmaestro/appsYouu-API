package branch

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type BranchQuery struct {
	Prisma *prisma.Client
}

func (r *BranchQuery) Branch(ctx context.Context, id string, language *string) (*prisma.Branch, error) {
	return r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)
}

func (r *BranchQuery) Branches(ctx context.Context, input gqlgen.BranchesInput, language *string) (*gqlgen.BranchConnection, error) {
	companyID := sessctx.CompanyWithFallback(ctx, r.Prisma, input.Company)

	branches, err := r.Prisma.Branches(&prisma.BranchesParams{
		Where: &prisma.BranchWhereInput{
			Company: &prisma.CompanyWhereInput{
				ID: &companyID,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	nodes := []*prisma.Branch{}

	for _, branch := range branches {
		clone := branch
		nodes = append(nodes, &clone)
	}

	return &gqlgen.BranchConnection{
		Nodes: nodes,
	}, nil
}

func New(client *prisma.Client) *BranchQuery {
	return &BranchQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
