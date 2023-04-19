package branch

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) Name(ctx context.Context, obj *prisma.Branch) (*string, error) {
	name, err := r.Prisma.Branch(prisma.BranchWhereUniqueInput{
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
