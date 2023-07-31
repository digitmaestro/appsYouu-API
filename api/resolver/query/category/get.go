package category

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *CategoryQuery) Category(ctx context.Context, id string) (*prisma.Category, error) {
	return r.Prisma.Category(prisma.CategoryWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)
}
