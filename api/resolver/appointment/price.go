package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/price"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) Price(ctx context.Context, obj *prisma.Appointment) (*gqlgen.Price, error) {
	return price.Convert(ctx, obj.Price), nil
}
