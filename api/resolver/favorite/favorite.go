package favorite

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type Favorite struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *Favorite {
	return &Favorite{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
