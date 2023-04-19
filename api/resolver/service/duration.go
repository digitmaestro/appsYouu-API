package service

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Service) Duration(ctx context.Context, obj *prisma.Service) (int, error) {
	if obj.Deleted {
		return 0, nil
	}

	return int(obj.Duration), nil
}
