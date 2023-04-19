package branch

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) News(ctx context.Context, obj *prisma.Branch) ([]*prisma.News, error) {
	newsItems, err := r.Prisma.Newses(&prisma.NewsesParams{
		Where: &prisma.NewsWhereInput{
			Branch: &prisma.BranchWhereInput{
				ID: &obj.ID,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	nodes := []*prisma.News{}

	for _, news := range newsItems {
		clone := news
		nodes = append(nodes, &clone)
	}

	return nodes, err
}
