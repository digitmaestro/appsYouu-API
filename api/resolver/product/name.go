package product

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Product) Name(ctx context.Context, obj *prisma.Product) (string, error) {
	if obj.Deleted {
		return i18n.Language(ctx)["DELETED_PRODUCT"], nil
	}

	name, err := r.Prisma.Product(prisma.ProductWhereUniqueInput{
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
