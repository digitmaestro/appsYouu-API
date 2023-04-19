package company

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type Company struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *Company {
	return &Company{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
