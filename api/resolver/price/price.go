package price

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
)

func Convert(ctx context.Context, price int32) *gqlgen.Price {
	v := float64(price) / 100
	return &gqlgen.Price{
		Value:        v,
		DisplayValue: i18n.FormatPrice(ctx, v),
	}
}
