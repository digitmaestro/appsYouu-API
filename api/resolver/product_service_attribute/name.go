package product_service_attribute

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *ProductServiceAttribute) Name(ctx context.Context, obj *prisma.ProductServiceAttribute) (*string, error) {
	name, err := r.Prisma.ProductServiceAttribute(prisma.ProductServiceAttributeWhereUniqueInput{
		ID: &obj.ID,
	}).Name().Exec(ctx)

	if err == prisma.ErrNoResult {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return i18n.GetLocalizedString(ctx, name), err
}
