package staff

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type StaffQuery struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *StaffQuery {
	return &StaffQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
