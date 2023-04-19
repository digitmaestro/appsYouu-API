package company

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

type CompanyQuery struct {
	Prisma *prisma.Client
}

func (r *CompanyQuery) Company(ctx context.Context, id string, language *string) (*prisma.Company, error) {
	return r.Prisma.Company(prisma.CompanyWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)
}

func (r *CompanyQuery) Companies(ctx context.Context, language *string) (*gqlgen.CompanyConnection, error) {
	companies, err := r.Prisma.Companies(&prisma.CompaniesParams{
		Where: &prisma.CompanyWhereInput{},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	nodes := []*prisma.Company{}

	for _, company := range companies {
		clone := company
		nodes = append(nodes, &clone)
	}

	return &gqlgen.CompanyConnection{
		Nodes: nodes,
	}, nil
}

func New(client *prisma.Client) *CompanyQuery {
	return &CompanyQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
