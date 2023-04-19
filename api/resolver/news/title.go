package news

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *News) Title(ctx context.Context, obj *prisma.News) (*string, error) {
	title, err := r.Prisma.News(prisma.NewsWhereUniqueInput{
		ID: &obj.ID,
	}).Title().Exec(ctx)

	if err == prisma.ErrNoResult {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return i18n.GetLocalizedString(ctx, title), err
}
