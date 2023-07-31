package category

import (
	"github.com/google/wire"
	"github.com/steebchen/keskin-api/prisma"
)

type Category struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *Category {
	return &Category{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
