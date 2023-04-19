package branch

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) Image(ctx context.Context, obj *prisma.Branch) (*gqlgen.Image, error) {
	return picture.FromID(obj.Image), nil
}
