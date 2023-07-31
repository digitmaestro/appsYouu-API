package category

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Category) SubCategories(ctx context.Context, obj *prisma.Category) ([]*prisma.SubCategory, error) {
	subCg, err := r.Prisma.SubCategories(&prisma.SubCategoriesParams{
		Where: &prisma.SubCategoryWhereInput{
			Category: &prisma.CategoryWhereInput{
				ID: &obj.ID,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, nil
	}

	nodes := []*prisma.SubCategory{}
	for _, cg := range subCg {
		clone := cg
		nodes = append(nodes, &clone)
	}

	return nodes, nil
}
