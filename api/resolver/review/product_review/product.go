package product_review

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *ProductReview) Product(ctx context.Context, obj *prisma.ProductReview) (*prisma.Product, error) {
	product, err := r.Prisma.Review(prisma.ReviewWhereUniqueInput{
		ID: &obj.ID,
	}).Product().Exec(ctx)

	return product, err
}
