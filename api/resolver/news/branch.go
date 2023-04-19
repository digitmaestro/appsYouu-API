package news

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *News) Branch(ctx context.Context, obj *prisma.News) (*prisma.Branch, error) {
	return r.Prisma.News(prisma.NewsWhereUniqueInput{
		ID: &obj.ID,
	}).Branch().Exec(ctx)
}
