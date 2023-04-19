package branch

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

var allowedTypes = []prisma.UserType{
	prisma.UserTypeManager,
}

type Branch struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *Branch {
	return &Branch{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
