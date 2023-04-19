package product

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/price"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Product) Price(ctx context.Context, obj *prisma.Product) (*gqlgen.Price, error) {
	if obj.Deleted {
		return price.Convert(ctx, 0), nil
	}

	return price.Convert(ctx, obj.Price), nil
}
