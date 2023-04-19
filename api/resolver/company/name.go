package company

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Company) Name(ctx context.Context, obj *prisma.Company) (*string, error) {
	name, err := r.Prisma.Company(prisma.CompanyWhereUniqueInput{
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
