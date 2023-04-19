package branch

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) Logo(ctx context.Context, obj *prisma.Branch) (*gqlgen.Image, error) {
	return picture.FromID(obj.Logo), nil
}
