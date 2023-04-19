package news

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/file"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *NewsMutation) CreateNews(
	ctx context.Context,
	input gqlgen.CreateNewsInput,
	language *string,
) (*gqlgen.CreateNewsPayload, error) {
	if err := permissions.CanAccessBranch(ctx, input.Branch, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	imageID, err := file.MaybeUpload(input.Data.Image, false)

	if err != nil {
		return nil, err
	}

	news, err := r.Prisma.CreateNews(prisma.NewsCreateInput{
		Title: *i18n.CreateLocalizedString(ctx, input.Data.Title),
		Image: imageID,

		Branch: prisma.BranchCreateOneWithoutNewsInput{
			Connect: &prisma.BranchWhereUniqueInput{
				ID: &input.Branch,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.CreateNewsPayload{
		News: news,
	}, nil
}
