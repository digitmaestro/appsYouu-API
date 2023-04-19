package company

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Company) Branches(ctx context.Context, obj *prisma.Company, pagination *gqlgen.PaginationInput) (*gqlgen.BranchConnection, error) {
	if pagination == nil {
		pagination = &gqlgen.PaginationInput{}
	}

	branches, err := r.Prisma.Branches(&prisma.BranchesParams{
		After: pagination.After,
		First: prisma.Int32Ptr(pagination.Limit),
		Where: &prisma.BranchWhereInput{
			Company: &prisma.CompanyWhereInput{
				ID: &obj.ID,
			},
		},
	}).Exec(ctx)

	nodes := []*prisma.Branch{}

	for _, branch := range branches {
		clone := branch
		nodes = append(nodes, &clone)
	}

	return &gqlgen.BranchConnection{
		Nodes: nodes,
	}, err
}
