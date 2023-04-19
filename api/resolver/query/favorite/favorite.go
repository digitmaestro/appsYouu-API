package favorite

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type FavoriteQuery struct {
	Prisma *prisma.Client
}

func (r *FavoriteQuery) Favorites(ctx context.Context, language *string) (*gqlgen.FavoritesPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	favorites, err := r.Prisma.Favorites(&prisma.FavoritesParams{
		Where: &prisma.FavoriteWhereInput{
			User: &prisma.UserWhereInput{
				ID: &viewer.ID,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	productsResult := []*prisma.Product{}
	uniqueProducts := make(map[string]bool)

	servicesResult := []*prisma.Service{}
	uniqueServices := make(map[string]bool)

	for _, favorite := range favorites {
		product, err := r.Prisma.Favorite(prisma.FavoriteWhereUniqueInput{
			ID: &favorite.ID,
		}).Product().Exec(ctx)

		if err != nil && err != prisma.ErrNoResult {
			return nil, err
		}

		if product != nil {
			_, known := uniqueProducts[product.ID]

			if !known {
				uniqueProducts[product.ID] = true
				clone := product
				productsResult = append(productsResult, clone)
			}
		}

		service, err := r.Prisma.Favorite(prisma.FavoriteWhereUniqueInput{
			ID: &favorite.ID,
		}).Service().Exec(ctx)

		if err != nil && err != prisma.ErrNoResult {
			return nil, err
		}

		if service != nil {
			_, known := uniqueServices[service.ID]

			if !known {
				uniqueServices[service.ID] = true
				clone := service
				servicesResult = append(servicesResult, clone)
			}
		}
	}

	return &gqlgen.FavoritesPayload{
		Products: &gqlgen.ProductConnection{
			Nodes: productsResult,
		},
		Services: &gqlgen.ServiceConnection{
			Nodes: servicesResult,
		},
	}, nil
}

func New(client *prisma.Client) *FavoriteQuery {
	return &FavoriteQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
