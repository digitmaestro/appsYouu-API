package company

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Company) Aliases(ctx context.Context, obj *prisma.Company) ([]string, error) {
	aliases, err := r.Prisma.Aliases(&prisma.AliasesParams{
		Where: &prisma.AliasWhereInput{
			Company: &prisma.CompanyWhereInput{
				ID: &obj.ID,
			},
		},
	}).Exec(ctx)

	nodes := []string{}

	for _, alias := range aliases {
		clone := alias.Value
		nodes = append(nodes, clone)
	}

	return nodes, err
}
