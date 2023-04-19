package favorite

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Favorite) Product(ctx context.Context, obj *prisma.Favorite) (*prisma.Product, error) {
	product, err := r.Prisma.Favorite(prisma.FavoriteWhereUniqueInput{
		ID: &obj.ID,
	}).Product().Exec(ctx)

	if err == prisma.ErrNoResult {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return product, err
}
