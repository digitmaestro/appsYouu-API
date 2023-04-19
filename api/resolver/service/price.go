package service

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/price"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Service) Price(ctx context.Context, obj *prisma.Service) (*gqlgen.Price, error) {
	if obj.Deleted {
		return price.Convert(ctx, 0), nil
	}

	return price.Convert(ctx, obj.Price), nil
}
