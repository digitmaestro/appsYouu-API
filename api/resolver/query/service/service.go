package service

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type ServiceQuery struct {
	Prisma *prisma.Client
}

func (r *ServiceQuery) Service(ctx context.Context, id string, language *string) (*prisma.Service, error) {
	return r.Prisma.Service(prisma.ServiceWhereUniqueInput{
		ID: &id,
	}).Exec(ctx)
}

func (r *ServiceQuery) Services(ctx context.Context, input gqlgen.ServiceInput, language *string) (*gqlgen.ServiceConnection, error) {
	deleted := false
	where := &prisma.ServiceWhereInput{
		Deleted: &deleted,
	}

	if input.GenderTarget != nil && *input.GenderTarget != prisma.GenderTargetAny {
		anyGender := prisma.GenderTargetAny
		where.Or = []prisma.ServiceWhereInput{
			{
				GenderTarget: input.GenderTarget,
			},
			{
				GenderTarget: &anyGender,
			},
		}
	}

	if input.Branch == nil {
		companyId := sessctx.CompanyWithFallback(ctx, r.Prisma, input.Company)

		where.Branch = &prisma.BranchWhereInput{
			Company: &prisma.CompanyWhereInput{
				ID: &companyId,
			},
		}
	} else {
		where.Branch = &prisma.BranchWhereInput{
			ID: input.Branch,
		}
	}

	services, err := r.Prisma.Services(&prisma.ServicesParams{
		Where:   where,
		OrderBy: assembleOrder(input.Order),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	nodes := []*prisma.Service{}

	for _, service := range services {
		clone := service
		nodes = append(nodes, &clone)
	}

	return &gqlgen.ServiceConnection{
		Nodes: nodes,
	}, nil
}

func New(client *prisma.Client) *ServiceQuery {
	return &ServiceQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
