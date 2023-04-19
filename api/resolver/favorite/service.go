package favorite

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Favorite) Service(ctx context.Context, obj *prisma.Favorite) (*prisma.Service, error) {
	service, err := r.Prisma.Favorite(prisma.FavoriteWhereUniqueInput{
		ID: &obj.ID,
	}).Service().Exec(ctx)

	if err == prisma.ErrNoResult {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return service, err
}
