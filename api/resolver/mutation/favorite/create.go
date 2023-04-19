package favorite

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *FavoriteMutation) CreateFavorite(
	ctx context.Context,
	input gqlgen.CreateFavoriteInput,
) (*gqlgen.CreateFavoritePayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if input.Data.Product == nil && input.Data.Service == nil {
		return nil, gqlerrors.NewValidationError("Product or service ID expected", "IDMissing")
	}

	if input.Data.Product != nil && input.Data.Service != nil {
		return nil, gqlerrors.NewValidationError("Either a product or service ID allowed", "IDMalformed")
	}

	query := prisma.FavoriteWhereInput{
		User: &prisma.UserWhereInput{
			ID: &viewer.ID,
		},
	}

	if input.Data.Product != nil {
		product, err := r.Prisma.Product(prisma.ProductWhereUniqueInput{
			ID: input.Data.Product,
		}).Exec(ctx)

		if err != nil || product.Deleted {
			return nil, gqlerrors.NewValidationError("Invalid product ID", "IDMalformed")
		}

		query.Product = &prisma.ProductWhereInput{
			ID: input.Data.Product,
		}
	}

	if input.Data.Service != nil {
		service, err := r.Prisma.Service(prisma.ServiceWhereUniqueInput{
			ID: input.Data.Service,
		}).Exec(ctx)

		if err != nil || service.Deleted {
			return nil, gqlerrors.NewValidationError("Invalid service ID", "IDMalformed")
		}

		query.Service = &prisma.ServiceWhereInput{
			ID: input.Data.Service,
		}
	}

	favorites, err := r.Prisma.Favorites(&prisma.FavoritesParams{
		Where: &query,
	}).Exec(ctx)

	if err != nil && err != prisma.ErrNoResult {
		return nil, err
	}

	if len(favorites) > 0 {
		return nil, gqlerrors.NewPermissionError("favorite already exists")
	}

	create := prisma.FavoriteCreateInput{
		User: prisma.UserCreateOneInput{
			Connect: &prisma.UserWhereUniqueInput{
				ID: &viewer.ID,
			},
		},
	}

	if input.Data.Product != nil {
		create.Product = &prisma.ProductCreateOneInput{
			Connect: &prisma.ProductWhereUniqueInput{
				ID: input.Data.Product,
			},
		}
	}

	if input.Data.Service != nil {
		create.Service = &prisma.ServiceCreateOneInput{
			Connect: &prisma.ServiceWhereUniqueInput{
				ID: input.Data.Service,
			},
		}
	}

	favorite, err := r.Prisma.CreateFavorite(create).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.CreateFavoritePayload{
		Favorite: favorite,
	}, nil
}
