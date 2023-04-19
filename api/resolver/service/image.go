package service

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Service) Image(ctx context.Context, obj *prisma.Service) (*gqlgen.Image, error) {
	if obj.Deleted {
		return nil, nil
	}

	return picture.FromID(obj.Image), nil
}
