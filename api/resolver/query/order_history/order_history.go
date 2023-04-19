package order_history

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type OrderHistoryQuery struct {
	Prisma *prisma.Client
}

func (r *OrderHistoryQuery) OrderHistory(ctx context.Context, language *string) (*gqlgen.OrderHistoryPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	productLinks, err := r.Prisma.AppointmentProductLinks(&prisma.AppointmentProductLinksParams{
		Where: &prisma.AppointmentProductLinkWhereInput{
			Appointment: &prisma.AppointmentWhereInput{
				Customer: &prisma.UserWhereInput{
					ID: &viewer.ID,
				},
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	productsResult := []*prisma.Product{}
	uniqueProducts := make(map[string]bool)

	for _, productLink := range productLinks {
		product, err := r.Prisma.AppointmentProductLink(prisma.AppointmentProductLinkWhereUniqueInput{
			ID: &productLink.ID,
		}).Product().Exec(ctx)

		if err != nil {
			return nil, err
		}

		_, known := uniqueProducts[product.ID]

		if !known {
			uniqueProducts[product.ID] = true
			clone := product
			productsResult = append(productsResult, clone)
		}
	}

	serviceLinks, err := r.Prisma.AppointmentServiceLinks(&prisma.AppointmentServiceLinksParams{
		Where: &prisma.AppointmentServiceLinkWhereInput{
			Appointment: &prisma.AppointmentWhereInput{
				Customer: &prisma.UserWhereInput{
					ID: &viewer.ID,
				},
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	servicesResult := []*prisma.Service{}
	uniqueServices := make(map[string]bool)

	for _, serviceLink := range serviceLinks {
		service, err := r.Prisma.AppointmentServiceLink(prisma.AppointmentServiceLinkWhereUniqueInput{
			ID: &serviceLink.ID,
		}).Service().Exec(ctx)

		if err != nil {
			return nil, err
		}

		_, known := uniqueServices[service.ID]

		if !known {
			uniqueServices[service.ID] = true
			clone := service
			servicesResult = append(servicesResult, clone)
		}
	}

	return &gqlgen.OrderHistoryPayload{
		Products: &gqlgen.ProductConnection{
			Nodes: productsResult,
		},
		Services: &gqlgen.ServiceConnection{
			Nodes: servicesResult,
		},
	}, nil
}

func New(client *prisma.Client) *OrderHistoryQuery {
	return &OrderHistoryQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
