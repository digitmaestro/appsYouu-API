package service

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Service) Desc(ctx context.Context, obj *prisma.Service) (*string, error) {
	if obj.Deleted {
		return nil, nil
	}

	desc, err := r.Prisma.Service(prisma.ServiceWhereUniqueInput{
		ID: &obj.ID,
	}).Desc().Exec(ctx)

	if err == prisma.ErrNoResult {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return i18n.GetLocalizedString(ctx, desc), err
}
