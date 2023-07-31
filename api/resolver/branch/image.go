package branch

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) Images(ctx context.Context, obj *prisma.Branch) ([]*gqlgen.Image, error) {
	return picture.FromIDsBranch(obj.Images)
}
