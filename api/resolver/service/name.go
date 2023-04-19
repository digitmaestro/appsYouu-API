package service

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Service) Name(ctx context.Context, obj *prisma.Service) (string, error) {
	if obj.Deleted {
		return i18n.Language(ctx)["DELETED_SERVICE"], nil
	}

	name, err := r.Prisma.Service(prisma.ServiceWhereUniqueInput{
		ID: &obj.ID,
	}).Name().Exec(ctx)

	if err == prisma.ErrNoResult {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	localizedValue := i18n.GetLocalizedString(ctx, name)

	if localizedValue == nil {
		return "", err
	} else {
		return *localizedValue, err
	}
}
