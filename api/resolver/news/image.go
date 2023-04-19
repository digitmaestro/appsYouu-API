package news

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *News) Image(ctx context.Context, obj *prisma.News) (*gqlgen.Image, error) {
	return picture.FromID(obj.Image), nil
}
