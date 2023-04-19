package appointment

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type Appointment struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *Appointment {
	return &Appointment{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
