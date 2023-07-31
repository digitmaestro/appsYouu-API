package sub_category

import (
	"github.com/google/wire"
	"github.com/steebchen/keskin-api/prisma"
)

type SubCategory struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *SubCategory {
	return &SubCategory{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
