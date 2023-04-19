package company

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Company) PwaIcon(ctx context.Context, obj *prisma.Company) (*gqlgen.Image, error) {
	return picture.FromID(obj.PwaIcon), nil
}
