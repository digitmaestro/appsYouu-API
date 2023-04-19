package product

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type Product struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *Product {
	return &Product{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
