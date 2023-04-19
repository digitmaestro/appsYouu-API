package product

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Product) Attributes(ctx context.Context, obj *prisma.Product) ([]*prisma.ProductServiceAttribute, error) {
	result := []*prisma.ProductServiceAttribute{}

	if !obj.Deleted {
		attributes, err := r.Prisma.Product(prisma.ProductWhereUniqueInput{
			ID: &obj.ID,
		}).Attributes(nil).Exec(ctx)

		if err != nil {
			return nil, err
		}

		for _, attribute := range attributes {
			clone := attribute
			result = append(result, &clone)
		}
	}

	return result, nil
}
