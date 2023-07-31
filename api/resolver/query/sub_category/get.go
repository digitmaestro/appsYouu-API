package sub_category

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *SubCategoryQuery) SubCategory(ctx context.Context, id string) (*prisma.SubCategory, error) {
	subCg, err := r.Prisma.SubCategory(prisma.SubCategoryWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return subCg, nil
}
