package company

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Company) PwaShortName(ctx context.Context, obj *prisma.Company) (string, error) {
	name, err := r.Prisma.Company(prisma.CompanyWhereUniqueInput{
		ID: &obj.ID,
	}).PwaShortName().Exec(ctx)

	if err == prisma.ErrNoResult {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	result := i18n.GetLocalizedString(ctx, name)

	if result != nil {
		return *result, err
	} else {
		return "", err
	}
}
