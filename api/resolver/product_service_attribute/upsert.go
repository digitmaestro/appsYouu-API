package product_service_attribute

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func UpsertAttributes(prismaClient *prisma.Client, ctx context.Context, productId *string, serviceId *string, attributes []*gqlgen.ProductServiceAttributeInput) error {
	if productId == nil && serviceId == nil {
		return nil
	}

	for _, attribute := range attributes {
		existingWhere := &prisma.ProductServiceAttributeWhereInput{
			Key: &attribute.Key,
		}

		if productId != nil {
			existingWhere.Product = &prisma.ProductWhereInput{
				ID: productId,
			}
		} else if serviceId != nil {
			existingWhere.Service = &prisma.ServiceWhereInput{
				ID: serviceId,
			}
		}

		existingAttributes, err := prismaClient.ProductServiceAttributes(&prisma.ProductServiceAttributesParams{
			Where: existingWhere,
		}).Exec(ctx)

		if err != nil {
			return err
		}

		if len(existingAttributes) == 0 {
			create := prisma.ProductServiceAttributeCreateInput{
				Key:   attribute.Key,
				Name:  *i18n.CreateLocalizedString(ctx, attribute.Name),
				Value: *i18n.CreateLocalizedString(ctx, attribute.Value),
			}

			if productId != nil {
				create.Product = &prisma.ProductCreateOneWithoutAttributesInput{
					Connect: &prisma.ProductWhereUniqueInput{
						ID: productId,
					},
				}
			} else if serviceId != nil {
				create.Service = &prisma.ServiceCreateOneWithoutAttributesInput{
					Connect: &prisma.ServiceWhereUniqueInput{
						ID: serviceId,
					},
				}
			}

			_, err = prismaClient.CreateProductServiceAttribute(create).Exec(ctx)
		} else {
			_, err = prismaClient.UpdateProductServiceAttribute(prisma.ProductServiceAttributeUpdateParams{
				Where: prisma.ProductServiceAttributeWhereUniqueInput{
					ID: &existingAttributes[0].ID,
				},
				Data: prisma.ProductServiceAttributeUpdateInput{
					Name:  i18n.UpdateRequiredLocalizedString(ctx, attribute.Name),
					Value: i18n.UpdateRequiredLocalizedString(ctx, attribute.Value),
				},
			}).Exec(ctx)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
