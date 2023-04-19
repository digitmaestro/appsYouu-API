package service

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Service) GenderTarget(ctx context.Context, obj *prisma.Service) (*prisma.GenderTarget, error) {
	if obj.Deleted {
		anyGender := prisma.GenderTargetAny
		return &anyGender, nil
	}

	return &obj.GenderTarget, nil
}
