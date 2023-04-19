package service

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Service) Attributes(ctx context.Context, obj *prisma.Service) ([]*prisma.ProductServiceAttribute, error) {
	result := []*prisma.ProductServiceAttribute{}

	if !obj.Deleted {
		attributes, err := r.Prisma.Service(prisma.ServiceWhereUniqueInput{
			ID: &obj.ID,
		}).Attributes(nil).Exec(ctx)

		if err != nil {
			return nil, err
		}

		for _, attribute := range attributes {
			clone := attribute
			result = append(result, &clone)
		}
	}

	return result, nil
}
