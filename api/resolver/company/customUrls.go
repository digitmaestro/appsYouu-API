package company

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Company) CustomUrls(ctx context.Context, obj *prisma.Company) ([]string, error) {
	customUrls, err := r.Prisma.CustomUrls(&prisma.CustomUrlsParams{
		Where: &prisma.CustomUrlWhereInput{
			Company: &prisma.CompanyWhereInput{
				ID: &obj.ID,
			},
		},
	}).Exec(ctx)

	nodes := []string{}

	for _, customUrl := range customUrls {
		clone := customUrl.Value
		nodes = append(nodes, clone)
	}

	return nodes, err
}
