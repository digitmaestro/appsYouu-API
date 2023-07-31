package sub_category

import (
	"github.com/google/wire"
	"github.com/steebchen/keskin-api/prisma"
)

type SubCategoryQuery struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *SubCategoryQuery {
	return &SubCategoryQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
