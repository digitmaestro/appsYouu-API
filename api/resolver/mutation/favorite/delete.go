package favorite

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *FavoriteMutation) DeleteFavorite(
	ctx context.Context,
	input gqlgen.DeleteFavoriteInput,
) (*gqlgen.DeleteFavoritePayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	deleted := false

	favorites, err := r.Prisma.Favorites(&prisma.FavoritesParams{
		Where: &prisma.FavoriteWhereInput{
			ID: &input.ID,
			User: &prisma.UserWhereInput{
				ID:      &viewer.ID,
				Deleted: &deleted,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(favorites) == 0 {
		return nil, gqlerrors.NewPermissionError("user has no permisson to access this favorite")
	}

	favorite, err := r.Prisma.DeleteFavorite(prisma.FavoriteWhereUniqueInput{
		ID: &input.ID,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.DeleteFavoritePayload{
		Favorite: favorite,
	}, nil
}
