package news

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *NewsMutation) DeleteNews(
	ctx context.Context,
	input gqlgen.DeleteNewsInput,
	language *string,
) (*gqlgen.DeleteNewsPayload, error) {
	branch, err := r.Prisma.News(prisma.NewsWhereUniqueInput{
		ID: &input.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	news, err := r.Prisma.DeleteNews(prisma.NewsWhereUniqueInput{
		ID: &input.ID,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.DeleteNewsPayload{
		News: news,
	}, nil
}
