package category

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *CategoryQuery) Categories(ctx context.Context) ([]*prisma.Category, error) {
	cgs, err := r.Prisma.Categories(&prisma.CategoriesParams{
		Where: &prisma.CategoryWhereInput{},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	nodes := []*prisma.Category{}

	for _, cg := range cgs {
		clone := cg
		nodes = append(nodes, &clone)
	}

	return nodes, nil
}
