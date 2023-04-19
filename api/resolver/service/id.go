package service

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Service) ID(ctx context.Context, obj *prisma.Service) (string, error) {
	if obj.Deleted {
		return "", nil
	}

	return obj.ID, nil
}
