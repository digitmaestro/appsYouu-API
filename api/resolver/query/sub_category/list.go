package sub_category

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *SubCategoryQuery) SubCategories(ctx context.Context) ([]*prisma.SubCategory, error) {
	cgs, err := r.Prisma.SubCategories(&prisma.SubCategoriesParams{
		Where: &prisma.SubCategoryWhereInput{},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	nodes := []*prisma.SubCategory{}

	for _, cg := range cgs {
		clone := cg
		nodes = append(nodes, &clone)
	}

	return nodes, nil
}
