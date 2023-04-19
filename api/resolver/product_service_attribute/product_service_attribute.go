package product_service_attribute

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type ProductServiceAttribute struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *ProductServiceAttribute {
	return &ProductServiceAttribute{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
