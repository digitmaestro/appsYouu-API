package appointment

import (
	"context"
	"errors"
	"time"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/prisma"
)

type PlanInput struct {
	Client          *prisma.Client
	Context         context.Context
	ProductRequests []*gqlgen.ConnectAppointmentProduct
	ServiceRequests []*gqlgen.ConnectAppointmentService
	Start           time.Time
}

type PlanPayload struct {
	// Price in cents
	Price int32

	// Duration in minutes
	Duration time.Duration

	Start string
	End   string

	CreateProducts []prisma.AppointmentProductLinkCreateWithoutAppointmentInput
	CreateServices []prisma.AppointmentServiceLinkCreateWithoutAppointmentInput
}

func Plan(input PlanInput) (*PlanPayload, error) {
	var cost int               // in cents
	var duration time.Duration // in minutes

	var productIDs []string

	for _, product := range input.ProductRequests {
		productIDs = append(productIDs, product.ID)
	}

	deleted := false

	products, err := input.Client.Products(&prisma.ProductsParams{
		Where: &prisma.ProductWhereInput{
			IDIn:    productIDs,
			Deleted: &deleted,
		},
	}).Exec(input.Context)

	if err != nil {
		return nil, err
	}

	var createProducts []prisma.AppointmentProductLinkCreateWithoutAppointmentInput

	for _, product := range input.ProductRequests {
		createProducts = append(createProducts, prisma.AppointmentProductLinkCreateWithoutAppointmentInput{
			Count: int32(product.Count),
			Product: prisma.ProductCreateOneInput{
				Connect: &prisma.ProductWhereUniqueInput{
					ID: &product.ID,
				},
			},
		})

		prodObj, err := findProduct(product.ID, products)

		if err != nil {
			return nil, err
		}

		for i := 0; i < product.Count; i++ {
			cost += int(prodObj.Price)
		}
	}

	var serviceIDs []string

	for _, service := range input.ServiceRequests {
		serviceIDs = append(serviceIDs, service.ID)
	}

	services, err := input.Client.Services(&prisma.ServicesParams{
		Where: &prisma.ServiceWhereInput{
			IDIn:    serviceIDs,
			Deleted: &deleted,
		},
	}).Exec(input.Context)

	if err != nil {
		return nil, err
	}

	var createServices []prisma.AppointmentServiceLinkCreateWithoutAppointmentInput

	for _, service := range input.ServiceRequests {
		createServices = append(createServices, prisma.AppointmentServiceLinkCreateWithoutAppointmentInput{
			Service: prisma.ServiceCreateOneInput{
				Connect: &prisma.ServiceWhereUniqueInput{
					ID: &service.ID,
				},
			},
		})

		svcObj, err := findService(service.ID, services)

		if err != nil {
			return nil, err
		}

		cost += int(svcObj.Price)
		duration += time.Duration(int(svcObj.Duration)) * time.Minute
	}

	start := *prisma.TimeString(input.Start)
	end := *prisma.TimeString(input.Start.Add(duration))

	return &PlanPayload{
		Price:          int32(cost),
		Duration:       duration,
		Start:          start,
		End:            end,
		CreateProducts: createProducts,
		CreateServices: createServices,
	}, nil
}

func findProduct(id string, products []prisma.Product) (prisma.Product, error) {
	for _, product := range products {
		if product.ID == id {
			return product, nil
		}
	}

	return prisma.Product{}, gqlerrors.NewFormatNodeError(errors.New(gqlerrors.PrismaNotFound), id)
}

func findService(id string, services []prisma.Service) (prisma.Service, error) {
	for _, service := range services {
		if service.ID == id {
			return service, nil
		}
	}

	return prisma.Service{}, gqlerrors.NewFormatNodeError(errors.New(gqlerrors.PrismaNotFound), id)
}
