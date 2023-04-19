package news

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type News struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *News {
	return &News{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
