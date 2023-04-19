package news

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/file"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *NewsMutation) UpdateNews(
	ctx context.Context,
	input gqlgen.UpdateNewsInput,
	language *string,
) (*gqlgen.UpdateNewsPayload, error) {
	branch, err := r.Prisma.News(prisma.NewsWhereUniqueInput{
		ID: &input.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessBranch(ctx, branch.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	imageID, err := file.MaybeUpload(input.Patch.Image, false)

	if err != nil {
		return nil, err
	}

	news, err := r.Prisma.UpdateNews(prisma.NewsUpdateParams{
		Where: prisma.NewsWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.NewsUpdateInput{
			Title: i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Title),
			Image: imageID,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateNewsPayload{
		News: news,
	}, nil
}
