package category

import (
	"github.com/google/wire"
	"github.com/steebchen/keskin-api/prisma"
)

type CategoryQuery struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *CategoryQuery {
	return &CategoryQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
